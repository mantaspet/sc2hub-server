package models

import (
	"database/sql"
	"github.com/mantaspet/sc2hub-server/pkg/validators"
)

type User struct {
	Username     string
	Password     string
	PasswordHash string
}

func (u User) Validate(db *sql.DB) map[string]string {
	errors := make(map[string]string)
	validators.SetError(errors, "Username",
		validators.Required(u.Username),
	)
	validators.SetError(errors, "Password",
		validators.Required(u.Password),
		validators.MinLength(u.Password, 8),
	)
	return validators.Errors(errors)
}
