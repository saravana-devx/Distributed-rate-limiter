package client

import (
	"context"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (c *Repository) CreateClient(ctx context.Context, client *Client) error {
	return c.db.WithContext(ctx).Create(client).Error

}

func (c *Repository) GetClientByID(ctx context.Context, clientID string) (*Client, error) {
	var client Client
	if err := c.db.WithContext(ctx).Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *Repository) GetClientByAPIKey(ctx context.Context, apiKey string) (*Client, error) {
	var client Client
	if err := c.db.WithContext(ctx).Where("api_key = ?", apiKey).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *Repository) GetAllClients(ctx context.Context) ([]Client, error) {
	var clients []Client
	if err := c.db.WithContext(ctx).Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func (c *Repository) UpdateClient(ctx context.Context, client *Client) error {
	return c.db.WithContext(ctx).Save(client).Error
}

func (c *Repository) DeleteClientByID(ctx context.Context, clientID string) error {
	return c.db.WithContext(ctx).Where("client_id = ?", clientID).Delete(&Client{}).Error
}
