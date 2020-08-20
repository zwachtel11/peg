package kube

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/component-base/version"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	Scheme = scheme.Scheme
)

type Client struct {
	kubeconfig string
}

func NewClient(kubeconfig string) Client {
	// If a kubeconfig file isn't provided, find one in the standard locations.
	if kubeconfig == "" {
		kubeconfig = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	}
	return Client{kubeconfig}
}

func (c *Client) getConfig() (*rest.Config, error) {
	config, err := clientcmd.LoadFromFile(c.kubeconfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load Kubeconfig file from %q", c.kubeconfig)
	}

	restConfig, err := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid configuration:") {
			return nil, errors.New(strings.Replace(err.Error(), "invalid configuration:", "invalid kubeconfig file", 1))
		}
		return nil, err
	}
	restConfig.UserAgent = fmt.Sprintf("peg/%s (%s)", version.Get().GitVersion, version.Get().Platform)

	// Set QPS and Burst to a threshold that ensures the controller runtime client/client go does't generate throttling log messages
	restConfig.QPS = 20
	restConfig.Burst = 100

	return restConfig, nil
}

func (c *Client) getKubeClient() (client.Client, error) {
	config, err := c.getConfig()
	if err != nil {
		return nil, err
	}

	var kc client.Client
	// Nb. The operation is wrapped in a retry loop to make newClientSet more resilient to temporary connection problems.
	connectBackoff := newConnectBackoff()
	if err := retryWithExponentialBackoff(connectBackoff, func() error {
		var err error
		kc, err = client.New(config, client.Options{Scheme: Scheme})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to connect to the management cluster")
	}

	return kc, nil
}

func (c *Client) Apply(objs []unstructured.Unstructured) error {
	createComponentObjectBackoff := newWriteBackoff()
	for i := range objs {
		obj := objs[i]

		// Create the Kubernetes object.
		// Nb. The operation is wrapped in a retry loop to make Create more resilient to unexpected conditions.
		if err := retryWithExponentialBackoff(createComponentObjectBackoff, func() error {
			return c.createObj(obj)
		}); err != nil {
			return err
		}
	}

	return nil
}

/**
variables, err := input.Processor.GetVariables(input.RawArtifact)
if err != nil {
	return nil, err
}

if input.ListVariablesOnly {
	return &template{
		variables:       variables,
		targetNamespace: input.TargetNamespace,
	}, nil
}

processedYaml, err := input.Processor.Process(input.RawArtifact, input.ConfigVariablesClient.Get)
if err != nil {
	return nil, err
}

// Transform the yaml in a list of objects, so following transformation can work on typed objects (instead of working on a string/slice of bytes).
objs, err := utilyaml.ToUnstructured(processedYaml)
if err != nil {
	return nil,
**/

func (c *Client) createObj(obj unstructured.Unstructured) error {
	kc, err := c.getKubeClient()
	if err != nil {
		return err
	}

	// check if the component already exists, and eventually update it
	currentR := &unstructured.Unstructured{}
	currentR.SetGroupVersionKind(obj.GroupVersionKind())

	key := client.ObjectKey{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}

	ctx := context.Background()
	if err := kc.Get(ctx, key, currentR); err != nil {
		if !apierrors.IsNotFound(err) {
			return errors.Wrapf(err, "failed to get current provider object")
		}

		//if it does not exists, create the component
		if err := kc.Create(ctx, &obj); err != nil {
			return errors.Wrapf(err, "failed to create provider object %s, %s/%s", obj.GroupVersionKind(), obj.GetNamespace(), obj.GetName())
		}
		return nil
	}

	obj.SetResourceVersion(currentR.GetResourceVersion())
	if err := kc.Patch(ctx, &obj, client.Merge); err != nil {
		return errors.Wrapf(err, "failed to patch provider object")
	}

	return nil
}
