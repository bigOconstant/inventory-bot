package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteI interface {
	Init() error
}

type Sqlite struct {
}

func (s *Sqlite) Init() (err error) {

	db, err := sql.Open("sqlite3", "./inventory.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS Settings (
		id INTEGER PRIMARY KEY,
		refresh_interval Integer NOT NULL DEFAULT 30,
		user_agent TEXT DEFAULT "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
		discord_webhook TEXT DEFAULT "",
		enabled BOOLEAN NOT NULL CHECK (enabled IN (0,1)) DEFAULT 1
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}
	return nil
}
