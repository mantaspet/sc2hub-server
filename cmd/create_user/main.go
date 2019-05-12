package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
)

var flgDsn = ""

func parseFlags() {
	flag.StringVar(&flgDsn, "dsn", "root:root@/sc2hub", "MySQL data source name")
	flag.Parse()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func createUser(db *sql.DB, username string, password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 symbols long")
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, password_hash) VALUES (?, ?)`

	_, err = db.Exec(stmt, username, passwordHash)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return errors.New("username already exists")
		}
		return err
	}

	return nil
}

func main() {
	parseFlags()

	var err error
	db, err := openDB(flgDsn + "?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')

	username = strings.TrimSuffix(username, "\n")
	password = strings.TrimSuffix(password, "\n")

	err = createUser(db, username, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("user " + username + " has been created")
}
