package client

import (
	"context"
	"ratelimiter/internal/redis"
	"time"
)

type ClientRepository interface {
	CreateClient(ctx context.Context, client *Client) error
	GetClientByID(ctx context.Context, clientID string) (*Client, error)
	GetClientByAPIKey(ctx context.Context, apiKey string) (*Client, error)
	GetAllClients(ctx context.Context) ([]Client, error)
	UpdateClient(ctx context.Context, client *Client) error
	DeleteClientByID(ctx context.Context, clientID string) error
}

type Service struct {
	repo  ClientRepository
	cache *redis.Redis
}

func NewService(repo ClientRepository, cache *redis.Redis) *Service {
	return &Service{repo: repo, cache: cache}
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
	s.cache.Set(ctx, data.ClientID, data, 5*time.Minute)
	return data, nil
}

func (s *Service) GetClientByIDService(ctx context.Context, clientID string) (*Client, error) {
	var cached Client
	if err := s.cache.Get(ctx, clientID, &cached); err == nil {
		return &cached, nil
	}

	client, err := s.repo.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	s.cache.Set(ctx, clientID, client, 5*time.Minute)
	return client, nil
}

func (s *Service) GetClientByAPIKeyService(ctx context.Context, apiKey string) (*Client, error) {
	var cached Client
	if err := s.cache.Get(ctx, apiKey, &cached); err == nil {
		return &cached, nil
	}

	client, err := s.repo.GetClientByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	s.cache.Set(ctx, apiKey, client, 5*time.Minute)
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
	s.cache.Del(ctx, existing.ClientID)
	s.cache.Del(ctx, existing.APIKey)
	return existing, nil
}

func (s *Service) DeleteClientByIDService(ctx context.Context, clientID string) error {
	existing, err := s.repo.GetClientByID(ctx, clientID)
	if err != nil {
		return err
	}
	if err := s.repo.DeleteClientByID(ctx, clientID); err != nil {
		return err
	}
	s.cache.Del(ctx, clientID)
	s.cache.Del(ctx, existing.APIKey)
	return nil
}
