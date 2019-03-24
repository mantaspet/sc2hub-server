package validators

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func SetError(errors map[string]string, field string, values ...string) {
	for _, val := range values {
		if len(val) != 0 {
			errors[field] = val
			return
		}
	}
}

func Errors(errors map[string]string) map[string]string {
	if len(errors) != 0 {
		return errors
	}
	return nil
}

func Required(value string) string {
	if strings.TrimSpace(value) == "" {
		return "Field is required"
	}
	return ""
}

func MaxLength(value string, d int) string {
	if utf8.RuneCountInString(value) > d {
		err := fmt.Sprintf("Field is too long (maximum is %d characters)", d)
		return err
	}
	return ""
}

func SQLUnique(db *sql.DB, table string, field string, value interface{}, id int) string {
	var res interface{}
	query := fmt.Sprintf("SELECT NULL FROM %s WHERE %s='%v'", table, field, value)
	if id > 0 {
		query += fmt.Sprintf(" AND id<>%v", id)
	}
	err := db.QueryRow(query).Scan(&res)
	if err != nil {
		return ""
	}
	return strings.Title(field) + " must be unique"
}

//func PermittedValues(field string, opts ...string) {
//    value := f.Get(field)
//    if value == "" {
//        return
//    }
//    for _, opt := range opts {
//        if value == opt {
//            return
//        }
//    }
//    f.Errors.Add(field, "This field is invalid")
//}

//func MinLength(field string, d int) {
//    value := f.Get(field)
//    if value == "" {
//        return
//    }
//    if utf8.RuneCountInString(value) < d {
//        f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
//    }
//}

//func MatchesPattern(field string, pattern *regexp.Regexp) {
//    value := f.Get(field)
//    if value == "" {
//        return
//    }
//    if !pattern.MatchString(value) {
//        f.Errors.Add(field, "This field is invalid")
//    }
//}
