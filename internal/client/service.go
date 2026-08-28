package client

import (
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateClientService(ctx context.Context, req *CreateClientRequest) (*Client, error) {
	existing, _ := s.repo.GetClientByID(ctx, req.ClientID)
	if existing != nil {
		return nil, ErrClientAlreadyExists
	}
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, err
	}
	data := &Client{
		ClientID:      req.ClientID,
		Algorithm:     req.Algorithm,
		Limit:         req.Limit,
		WindowSeconds: req.WindowSeconds,
		APIKey:        apiKey,
	}
	err = s.repo.CreateClient(ctx, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetClientByIDService(ctx context.Context, clientID string) (*Client, error) {
	client, err := s.repo.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *Service) GetClientByAPIKeyService(ctx context.Context, apiKey string) (*Client, error) {
	client, err := s.repo.GetClientByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *Service) GetAllClientsService(ctx context.Context) ([]Client, error) {
	clients, err := s.repo.GetAllClients(ctx)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (s *Service) UpdateClientService(ctx context.Context, clientID string, req *UpdateClientRequest) (*Client, error) {
	existing, err := s.repo.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	existing.Algorithm = req.Algorithm
	existing.Limit = req.Limit
	existing.WindowSeconds = req.WindowSeconds
	err = s.repo.UpdateClient(ctx, existing)
	if err != nil {
		return nil, err
	}
	return existing, nil

}

func (s *Service) DeleteClientByIDService(ctx context.Context, clientID string) error {
	err := s.repo.DeleteClientByID(ctx, clientID)
	if err != nil {
		return err
	}
	return nil
}
