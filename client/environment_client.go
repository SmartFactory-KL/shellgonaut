package client

import (
	"fmt"
	"log/slog"
)

type AASEnvironmentClient struct {
	ShellRepositoryClient
	SubmodelRepositoryClient
}

type AASEnvironmentClientConfig struct {
	ShellRepositoryBaseUrl    string
	SubmodelRepositoryBaseUrl string
}

// NewAASEnvironmentClient creates a ShellClient with integrated http client
func NewAASEnvironmentClient(cfg *AASEnvironmentClientConfig) (*AASEnvironmentClient, error) {
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

	return &AASEnvironmentClient{
		ShellRepositoryClient:    *shellRepoClient,
		SubmodelRepositoryClient: *submodelRepoClient,
	}, nil
}

// CheckTargetAvailability calls the /description endpoint of both repos
// as a basic health check
func (envClient *AASEnvironmentClient) CheckTargetAvailability() error {
	if _, err := envClient.GetShellRepositoryDescription(); err != nil {
		return fmt.Errorf("shell repository description not readable: %w", err)
	}

	if _, err := envClient.GetSubmodelRepositoryDescription(); err != nil {
		return fmt.Errorf("submodel repository description not readable: %w", err)
	}

	slog.Info("BasyxRepositoryClient: Repositories Check OK")

	return nil
}
