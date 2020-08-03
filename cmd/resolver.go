package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	auth "github.com/zwachtel11/peg/pkg/auth/docker"
)

func newResolver(username, password string) remotes.Resolver {
	if username != "" || password != "" {
		return docker.NewResolver(docker.ResolverOptions{
			Credentials: func(hostName string) (string, string, error) {
				return username, password, nil
			},
		})
	}

	cli, err := auth.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading auth file: %v\n", err)
	}

	resolver, err := cli.Resolver(context.Background(), http.DefaultClient, false)
	if err != nil {
		resolver = docker.NewResolver(docker.ResolverOptions{})
	}

	return resolver
}
