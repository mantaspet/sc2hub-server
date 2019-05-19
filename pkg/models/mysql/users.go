package mysql

import (
	"database/sql"
	"github.com/mantaspet/sc2hub-server/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) SelectOne(username string) (*models.User, error) {
	stmt := `SELECT username, password_hash FROM users WHERE username=?`

	user := &models.User{}
	err := m.DB.QueryRow(stmt, username).Scan(&user.Username, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}
