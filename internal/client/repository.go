package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const pgUniqueViolationCode = "23505"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (c *Repository) CreateClient(ctx context.Context, client *Client) error {
	if err := c.db.WithContext(ctx).Create(client).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return ErrClientAlreadyExists
		}
		return fmt.Errorf("CreateClient: %w", err)
	}
	return nil
}

func (c *Repository) GetClientByID(ctx context.Context, clientID string) (*Client, error) {
	var client Client
	if err := c.db.WithContext(ctx).Where("client_id = ?", clientID).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, fmt.Errorf("GetClientByID: %w", err)
	}
	return &client, nil
}

func (c *Repository) GetClientByAPIKey(ctx context.Context, apiKey string) (*Client, error) {
	var client Client
	if err := c.db.WithContext(ctx).Where("api_key = ?", apiKey).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, fmt.Errorf("GetClientByAPIKey: %w", err)
	}
	return &client, nil
}

func (c *Repository) GetAllClients(ctx context.Context) ([]Client, error) {
	var clients []Client
	if err := c.db.WithContext(ctx).Find(&clients).Error; err != nil {
		return nil, fmt.Errorf("GetAllClients: %w", err)
	}
	return clients, nil
}

func (c *Repository) UpdateClient(ctx context.Context, client *Client) error {
	if err := c.db.WithContext(ctx).Save(client).Error; err != nil {
		return fmt.Errorf("UpdateClient: %w", err)
	}
	return nil
}

func (c *Repository) DeleteClientByID(ctx context.Context, clientID string) error {
	if err := c.db.WithContext(ctx).Where("client_id = ?", clientID).Delete(&Client{}).Error; err != nil {
		return fmt.Errorf("DeleteClientByID: %w", err)
	}
	return nil
}
