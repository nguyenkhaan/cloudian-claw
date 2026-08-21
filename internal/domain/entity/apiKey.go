package entity

import (
	"errors"
	"time"
)

type APIKey struct {
	ID      string `json:"id"`
	Revoked bool   `json:"revoked"`
	UserID  string `json:"user_id"`
}

type APIKeyUsage struct {
	Key      APIKey    `json:"api_key"`
	Endpoint string    `json:"endpoint"`
	Status   string    `json:"status,omitempty"`
	Time     time.Time `json:"time"`
}

func (k APIKey) Validate() error {
	if k.ID == "" {
		return errors.New("validate api_key: id is required")
	}
	if k.UserID == "" {
		return errors.New("validate api_key: user_id is required")
	}
	return nil
}

func (u APIKeyUsage) Validate() error {
	if u.Key.ID == "" {
		return errors.New("validate api_key_usage: key.id is required")
	}
	if u.Endpoint == "" {
		return errors.New("validate api_key_usage: endpoint is required")
	}
	return nil
}
