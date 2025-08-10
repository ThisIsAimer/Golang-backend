package models

type Execs struct {
	ID        int    `json:"id,omitempty" db:"id,omitempty"`
	FirstName string `json:"first_name,omitempty" db:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty" db:"last_name,omitempty"`
}
