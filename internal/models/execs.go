package models

import "database/sql"

type Execs struct {
	ID        int    `json:"id,omitempty" db:"id,omitempty"`
	FirstName string `json:"first_name,omitempty" db:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty" db:"last_name,omitempty"`
	Email     string `json:"email,omitempty" db:"email,omitempty"`
	Role      string `json:"role,omitempty" db:"role,omitempty"`
	UserName  string `json:"user_name,omitempty" db:"user_name,omitempty"`
	Password  string `json:"password,omitempty" db:"password,omitempty"`

	// sql.NullStrings will automatically updated by the database so we dont need to populate it
	// it can be null or a string
	UserCreatedAt     sql.NullString `json:"user_created_at,omitempty" db:"user_created_at,omitempty"`
	PasswordChangedAt sql.NullString `json:"password_changed_at,omitempty" db:"password_changed_at,omitempty"`
	PassResetCode     sql.NullString `json:"pass_reset_code,omitempty" db:"pass_reset_code,omitempty"`
	PassCodeExpires   sql.NullString `json:"pass_code_expires,omitempty" db:"pass_code_expires,omitempty"`

	InactiveStatus bool `json:"inactive_status,omitempty" db:"inactive_status,omitempty"`
}

// for normal use
type BasicExecs struct {
	ID             int    `json:"id,omitempty" db:"id,omitempty"`
	UserName       string `json:"user_name,omitempty" db:"user_name,omitempty"`
	Password       string `json:"password,omitempty" db:"password,omitempty"`
	FirstName      string `json:"first_name,omitempty" db:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty" db:"last_name,omitempty"`
	Email          string `json:"email,omitempty" db:"email,omitempty"`
	UserCreatedAt  string `json:"user_created_at,omitempty" db:"user_created_at,omitempty"`
	Role           string `json:"role,omitempty" db:"role,omitempty"`
	InactiveStatus bool   `json:"inactive_status,omitempty" db:"inactive_status,omitempty"`
}
