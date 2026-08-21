package entity

type User struct {
	ID          string
	Name        string `json:"name,omitempty"`
	Information string `json:"information,omitempty"`
}
