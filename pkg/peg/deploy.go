package peg

import (
	"context"

	"github.com/containerd/containerd/remotes"
	"github.com/zwachtel11/peg/pkg/kube"
)

// Deploy ...
func Deploy(ctx context.Context, resolver remotes.Resolver, manifest string, kubeconfig string) error {
	bytes, err := Pull(ctx, resolver, manifest)
	if err != nil {
		return err
	}

	unstruct, err := kube.PrepareManifest(bytes)
	if err != nil {
		return err
	}

	cli := kube.NewClient(kubeconfig)
	err = cli.Apply(unstruct)
	if err != nil {
		return err
	}

	return nil
}
