package kube

import (
	"bufio"
	"bytes"
	"io"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	apiyaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

const namespaceKind = "Namespace"

func PrepareManifest(bytes []byte) ([]unstructured.Unstructured, error) {
	unstruct, err := ToUnstructured(bytes)
	if err != nil {
		return nil, err
	}
	return fixTargetNamespace(unstruct, "default"), nil
}

// fixTargetNamespace ensures all the provider components are deployed in the target namespace (apply only to namespaced objects).
func fixTargetNamespace(objs []unstructured.Unstructured, targetNamespace string) []unstructured.Unstructured {
	for _, o := range objs {
		// if the object has Kind Namespace, fix the namespace name
		if o.GetKind() == namespaceKind {
			o.SetName(targetNamespace)
		}

		// if the object is namespaced, set the namespace name
		if IsResourceNamespaced(o.GetKind()) {
			o.SetNamespace(targetNamespace)
		}
	}

	return objs
}

// IsResourceNamespaced returns true if the resource kind is namespaced.
func IsResourceNamespaced(kind string) bool {
	switch kind {
	case "Namespace",
		"Node",
		"PersistentVolume",
		"PodSecurityPolicy",
		"CertificateSigningRequest",
		"ClusterRoleBinding",
		"ClusterRole",
		"VolumeAttachment",
		"StorageClass",
		"CSIDriver",
		"CSINode",
		"ValidatingWebhookConfiguration",
		"MutatingWebhookConfiguration",
		"CustomResourceDefinition",
		"PriorityClass",
		"RuntimeClass":
		return false
	default:
		return true
	}
}

// ToUnstructured takes a YAML and converts it to a list of Unstructured objects
func ToUnstructured(rawyaml []byte) ([]unstructured.Unstructured, error) {
	var ret []unstructured.Unstructured

	reader := apiyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(rawyaml)))
	count := 1
	for {
		// Read one YAML document at a time, until io.EOF is returned
		b, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrapf(err, "failed to read yaml")
		}
		if len(b) == 0 {
			break
		}

		var m map[string]interface{}
		if err := yaml.Unmarshal(b, &m); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal the yaml document")
		}
		var u unstructured.Unstructured
		u.SetUnstructuredContent(m)

		// Ignore empty objects.
		// Empty objects are generated if there are weird things in manifest files like e.g. two --- in a row without a yaml doc in the middle
		if u.Object == nil {
			continue
		}

		ret = append(ret, u)
		count++
	}

	return ret, nil
}

// retryWithExponentialBackoff repeats an operation until it passes or the exponential backoff times out.
func retryWithExponentialBackoff(opts wait.Backoff, operation func() error) error {

	i := 0
	err := wait.ExponentialBackoff(opts, func() (bool, error) {
		i++
		if err := operation(); err != nil {
			if i < opts.Steps {
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	if err != nil {
		return errors.Wrapf(err, "action failed after %d attempts", i)
	}
	return nil
}

// newWriteBackoff creates a new API Machinery backoff parameter set suitable for use with clusterctl write operations.
func newWriteBackoff() wait.Backoff {
	// Return a exponential backoff configuration which returns durations for a total time of ~40s.
	// Example: 0, .5s, 1.2s, 2.3s, 4s, 6s, 10s, 16s, 24s, 37s
	// Jitter is added as a random fraction of the duration multiplied by the jitter factor.
	return wait.Backoff{
		Duration: 500 * time.Millisecond,
		Factor:   1.5,
		Steps:    5,
		Jitter:   0.4,
	}
}

// newConnectBackoff creates a new API Machinery backoff parameter set suitable for use when clusterctl connect to a cluster.
func newConnectBackoff() wait.Backoff {
	// Return a exponential backoff configuration which returns durations for a total time of ~15s.
	// Example: 0, .25s, .6s, 1.2, 2.1s, 3.4s, 5.5s, 8s, 12s
	// Jitter is added as a random fraction of the duration multiplied by the jitter factor.
	return wait.Backoff{
		Duration: 250 * time.Millisecond,
		Factor:   1.5,
		Steps:    9,
		Jitter:   0.1,
	}
}

// newReadBackoff creates a new API Machinery backoff parameter set suitable for use with clusterctl read operations.
func newReadBackoff() wait.Backoff {
	// Return a exponential backoff configuration which returns durations for a total time of ~15s.
	// Example: 0, .25s, .6s, 1.2, 2.1s, 3.4s, 5.5s, 8s, 12s
	// Jitter is added as a random fraction of the duration multiplied by the jitter factor.
	return wait.Backoff{
		Duration: 250 * time.Millisecond,
		Factor:   1.5,
		Steps:    9,
		Jitter:   0.1,
	}
}
