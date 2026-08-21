package entity

import (
	"errors"
)

type User struct {
	ID          string
	Name        string `json:"name,omitempty"`
	Information string `json:"information,omitempty"`
}

func (u User) Validate() error {
	if u.ID == "" {
		return errors.New("validate user: id is required")
	}
	return nil
}
