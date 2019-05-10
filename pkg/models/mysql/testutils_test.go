package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mantaspet/sc2hub-server/pkg/models/mock"
	"io/ioutil"
	"testing"
)

func seedData(db *sql.DB) error {
	stmt := ""
	stmt += "INSERT INTO platforms (id, name) VALUES (1, 'twitch'), (2, 'youtube');"

	stmt += ""
	for _, ec := range mock.EventCategories {
		stmt += fmt.Sprintf("INSERT INTO event_categories"+
			"(name, pattern, info_url, image_url, description, priority) "+
			"VALUES ('%v', '%v', '%v', '%v', '%v', %v);",
			ec.Name, ec.Pattern, ec.InfoURL, ec.ImageURL, ec.Description, ec.Priority)
	}

	for _, c := range mock.Channels {
		stmt += fmt.Sprintf("INSERT INTO channels (id, platform_id, login, title, profile_image_url) "+
			"VALUES ('%v', %v, '%v', '%v', '%v');", c.ID, c.PlatformID, c.Login, c.Title, c.ProfileImageURL)
	}

	stmt += `INSERT INTO event_category_channels (event_category_id, channel_id)
		VALUES (1, '42508152'), (2, 'UCK5eBtuoj_HkdXKHNmBLAXg');`

	fmt.Println(stmt)
	_, err := db.Exec(stmt)
	return err
}

func newTestDB(t *testing.T) (*sql.DB, func()) {
	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contains multiple SQL statements, we need
	// to use the `multiStatements=true` parameter in our DSN. This instructs
	// our MySQL database driver to support executing multiple SQL statements
	// in one db.Exec()` call.
	db, err := sql.Open("mysql", "test_web:test@/test_sc2hub?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script from file and execute the statements.
	script, err := ioutil.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	err = seedData(db)
	if err != nil {
		t.Fatal(err)
	}

	// Return the connection pool and an anonymous function which reads and
	// executes the teardown script, and closes the connection pool. We can
	// assign this anonymous function and call it later once our test has
	// completed.
	return db, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		_ = db.Close()
	}
}
