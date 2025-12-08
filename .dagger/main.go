package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"dagger/blog-ci/internal/dagger"
)

type BlogCi struct{}

// Build container from Dockerfile
func (m *BlogCi) BuildFromDockerfile(
	// +defaultPath="/"
	source *dagger.Directory) *dagger.Container {
	return dag.Container().Build(source)
}


// Publish Docker image to registry
func (m *BlogCi) PublishImage(ctx context.Context, name string,
	// +default="latest"
	tag string,
	// +default="ttl.sh"
	registry string,
	username string,
	password *dagger.Secret,
	// +defaultPath="/"
	source *dagger.Directory,
) (string, error) {

	container := m.BuildFromDockerfile(source)	

	if registry != "ttl.sh" {
		container.WithRegistryAuth(registry, username, password)
		return container.Publish(ctx, fmt.Sprintf("%s/%s/%s:%s", registry, username, name, tag))
	} else {
		return container.Publish(ctx, fmt.Sprintf("%s/%s-%.0f", registry, name, math.Floor(rand.Float64()*10000000)))
	}
}
