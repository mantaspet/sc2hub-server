package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func fieldExists(table string, field string, value interface{}) error {
	var res interface{}
	query := fmt.Sprintf("SELECT NULL FROM %s WHERE %s='%v'", table, field, value)
	err := db.QueryRow(query).Scan(&res)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return nil
	}
	return errors.New(strings.Title(field) + " must be unique")
}
