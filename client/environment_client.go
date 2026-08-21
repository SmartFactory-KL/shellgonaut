package client

import (
	"fmt"
	"log/slog"
)

type AASEnvironmentClient struct {
	ShellRepositoryClient
	SubmodelRepositoryClient
	ConceptDescriptionRepositoryClient
}

type AASEnvironmentClientConfig struct {
	ShellRepositoryBaseUrl    string
	SubmodelRepositoryBaseUrl string
	ConceptDescriptionBaseUrl string
}

// NewAASEnvironmentClient creates a ShellClient with integrated http client
func NewAASEnvironmentClient(cfg *AASEnvironmentClientConfig, opts ...ClientOption) (*AASEnvironmentClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("invalid config: cannot be nil")
	}

	shellRepoClient, err := NewShellRepositoryClient(cfg.ShellRepositoryBaseUrl, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create shell repository client: %w", err)
	}

	submodelRepoClient, err := NewSubmodelRepositoryClient(cfg.SubmodelRepositoryBaseUrl, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create submodel repository client: %w", err)
	}

	cdRepoClient, err := NewConceptDescriptionRepositoryClient(cfg.ConceptDescriptionBaseUrl, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create concept description repository client: %w", err)
	}

	return &AASEnvironmentClient{
		ShellRepositoryClient:              *shellRepoClient,
		SubmodelRepositoryClient:           *submodelRepoClient,
		ConceptDescriptionRepositoryClient: *cdRepoClient,
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

	if _, err := envClient.GetConceptDescriptionRepositoryDescription(); err != nil {
		return fmt.Errorf("concept description repository description not readable: %w", err)
	}

	slog.Info("AASEnvironmentClient: Repositories Check OK")

	return nil
}
