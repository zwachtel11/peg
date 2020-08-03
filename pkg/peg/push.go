package peg

import (
	"context"
	"io/ioutil"

	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// Push
func Push(ctx context.Context, resolver remotes.Resolver, manifest, filepath string) error {
	fileContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Push file(s) w custom mediatype to registry
	memoryStore := content.NewMemoryStore()
	desc := memoryStore.Add(filepath, yamlMediaType, fileContent)
	pushContents := []ocispec.Descriptor{desc}

	desc, err = oras.Push(ctx, resolver, manifest, memoryStore, pushContents)
	if err != nil {
		return errors.Wrapf(err, "oras pull failed %s", manifest)
	}

	return nil
}
