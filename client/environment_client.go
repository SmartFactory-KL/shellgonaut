package client

import (
	"fmt"
	"log/slog"
)

type BasyxEnvironmentClient struct {
	ShellRepositoryClient
	SubmodelRepositoryClient
}

type BasyxRepositoryClientConfig struct {
	ShellRepositoryBaseUrl    string
	SubmodelRepositoryBaseUrl string
}

// TODO: This module is intended to become a standalone Go based Basyx client for internal use
// TODO: to be published on GitHub. Thats why it is so extensive for now.

// TODO: General Notes on how this can become a Standalone Client
// TODO: For only a single repository, this would be fine. But since registries are a thing,
// TODO: one would need to create a single BasyxRepoClient for all repositories present within a registry
// TODO: whenever they come up. So this RepositoryClient should be accessible through another layer,
// TODO: which could be a "SingleRepoClient" or a "RegistryClient" or even a "DiscoveryClient"

// NewBasyxRepositoryClient creates a ShellClient with integrated http client
func NewBasyxRepositoryClient(cfg *BasyxRepositoryClientConfig) (*BasyxEnvironmentClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("invalid config: cannot be nil")
	}

	shellRepoClient, err := NewShellRepositoryClient(cfg.ShellRepositoryBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create shell repository client: %w", err)
	}

	submodelRepoClient, err := NewSubmodelRepositoryClient(cfg.SubmodelRepositoryBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create submodel repository client: %w", err)
	}

	return &BasyxEnvironmentClient{
		ShellRepositoryClient:    *shellRepoClient,
		SubmodelRepositoryClient: *submodelRepoClient,
	}, nil
}

// CheckTargetAvailability calls the /description endpoint of both repos
// as a basic health check
func (envClient *BasyxEnvironmentClient) CheckTargetAvailability() error {
	if _, err := envClient.GetShellRepositoryDescription(); err != nil {
		return fmt.Errorf("shell repository description not readable: %w", err)
	}

	if _, err := envClient.GetSubmodelRepositoryDescription(); err != nil {
		return fmt.Errorf("submodel repository description not readable: %w", err)
	}

	slog.Info("BasyxRepositoryClient: Repositories Check OK")

	return nil
}
