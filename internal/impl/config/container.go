package config

import (
	"context"
	"fmt"
)

// Khoi tao cac thanh phan cua ung dung
type Container struct {
	Config *ApplicationConfig
}

func NewContainer(ctx context.Context) (*Container, error) {
	// Cac ctx se thcu hien vai tro sau 
	applicationConfig, err := LoadApplicationConfig() 
	if err != nil {
		return nil, fmt.Errorf("Loading application config error: ", err)
	} 


	return &Container{
		Config : applicationConfig, 
	}, nil 
}

func (c *Container) Close() error {
	// Tien hanh cac buoc dong ket noi nhu db connection, http connection.
	return nil
}
