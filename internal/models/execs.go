package models

import "database/sql"

type Execs struct {
	ID        int    `json:"id,omitempty" db:"id,omitempty"`
	FirstName string `json:"first_name,omitempty" db:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty" db:"last_name,omitempty"`
	Email     string `json:"email,omitempty" db:"email,omitempty"`
	Role      string
	UserName  string
	Password  string
	
	// sql.NullStrings will automatically updated by the database so we dont need to populate it
	// it can be null or a string
	UserCreatedAt     sql.NullString
	PasswordChangedAt sql.NullString
	PassResetCode     sql.NullString
	PassCodeExpires   sql.NullString

	InactiveStatus    bool
}
