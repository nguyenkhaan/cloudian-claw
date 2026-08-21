package entity

import "time"

type APIKey struct {
	ID string `json:"id"`
	Revoked bool `json:"revoked"`
	UserID string `json:"user_id"`
}

type APIKeyUsage struct {
	Key APIKey `json:"api_key"`
	Endpoint string `json:"endpoint"`
	Status string `json:"status,omitempty"`  //failed, success 
	Time time.Time `json:"time"`
}