package sqlite

import (
	"database/sql"
	"fmt"
	"goinventory/db/dbmodels"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteI interface {
	Init() error
	SaveSettings(settings dbmodels.Settings) error
	GetSettings() (settings dbmodels.Settings, err error)
}

type Sqlite struct {
}

func (s *Sqlite) Init() (err error) {

	db, err := sql.Open("sqlite3", "./inventory.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/* Create and init settings table */
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

	sqlStmt = `
		select count(*) from Settings;
	`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for rows.Next() {

		err = rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Count :", count)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()

	if count == 0 {
		fmt.Println("First run creating settings object")
		sqlStmt = "insert into Settings (id) Values (1);"
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return err
		}
	}

	/* Create and init Items table table */

	sqlStmt = `
	CREATE TABLE IF NOT EXISTS Items (
		id INTEGER PRIMARY KEY,
		url TEXT DEFAULT "",
		item TEXT DEFAULT ""
	);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	return nil
}

func (s *Sqlite) GetSettings() (settings dbmodels.Settings, err error) {
	db, err := sql.Open("sqlite3", "./inventory.db")
	if err != nil {
		fmt.Println("error preparing")
		log.Fatal(err)
	}
	defer db.Close()
	query := `
			select refresh_interval,
			user_agent,
			discord_webhook,
			enabled from settings;
	`
	rows, err := db.Query(query)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&settings.Refresh_interval, &settings.User_agent, &settings.Discord_webhook, &settings.Enabled)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(settings)
	}
	return
}

func (s *Sqlite) SaveSettings(settings dbmodels.Settings) (err error) {
	db, err := sql.Open("sqlite3", "./inventory.db")
	if err != nil {
		fmt.Println("error preparing")
		log.Fatal(err)
	}
	defer db.Close()
	/*
		id,refresh_interval , user_agent , discord_webhook , enabled
	*/
	updatestmt := `
			update Settings set refresh_interval = ?,
			user_agent = ?,
			discord_webhook = ?,
			enabled = ?;
	`
	stmt, err := db.Prepare(updatestmt)
	defer stmt.Close()
	fmt.Println("done preparing")
	fmt.Println(settings)
	updates := []interface{}{
		settings.Refresh_interval,
		settings.User_agent,
		settings.Discord_webhook,
		settings.Enabled,
	}

	_, err = stmt.Exec(updates...)

	if err != nil {
		fmt.Println("error updating settings:", err)
	}
	fmt.Println("returning")
	return err
}
