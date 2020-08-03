package peg

import (
	"context"

	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	"github.com/pkg/errors"
)

const yamlMediaType = "application/vnd.ecp.layer.v1+yaml"

// Pull ...
func Pull(ctx context.Context, resolver remotes.Resolver, manifest string) ([]byte, error) {
	memStore := content.NewMemoryStore()
	allowedMediaTypes := []string{yamlMediaType}
	_, content, err := oras.Pull(ctx, resolver, manifest, memStore, oras.WithAllowedMediaTypes(allowedMediaTypes))
	if err != nil {
		return nil, errors.Wrapf(err, "oras pull failed %s", manifest)
	}

	if len(content) != 1 {
		return nil, errors.Errorf("%s has invalid configuration", manifest)
	}

	_, bytes, ok := memStore.Get(content[0])
	if !ok {
		return nil, err
	}

	return bytes, nil
}
