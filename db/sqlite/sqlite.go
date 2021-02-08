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
	Close()
	SaveSettings(settings dbmodels.Settings) error
	SaveItem(name string, url string) (id int, err error)
	GetItems() (items []dbmodels.Item, err error)
	GetSettings() (settings dbmodels.Settings, err error)
	DeleteItem(id int) error
}

type Sqlite struct {
	db *sql.DB
}

func (s *Sqlite) Close() {
	fmt.Println("calling close")
	s.db.Close()
}

func (s *Sqlite) Init() (err error) {

	s.db, err = sql.Open("sqlite3", "./inventory.db")
	if err != nil {
		log.Fatal(err)
	}

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
	_, err = s.db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	sqlStmt = `
		select count(*) from Settings;
	`
	rows, err := s.db.Query(sqlStmt)
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
		_, err = s.db.Exec(sqlStmt)
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

	_, err = s.db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	return nil
}

func (s *Sqlite) DeleteItem(id int) error {
	sqlStmt := `
		delete from Items where id = ?;
	`
	stmt, err := s.db.Prepare(sqlStmt)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(&id)

	return err
}
func (s *Sqlite) SaveItem(name string, url string) (id int, err error) {

	if err != nil {
		fmt.Println("error opening db")
		log.Printf("%q\n", err)
		return
	}

	insertstmt := `
			INSERT INTO Items (url,item) VALUES (?,?)
	`
	stmt, err := s.db.Prepare(insertstmt)
	defer stmt.Close()
	inserted, err := stmt.Exec(&url, &name)
	if err != nil {
		fmt.Println("error inserting")
		log.Printf("%q\n", err)
		return
	}

	id64, err := inserted.LastInsertId()
	id = int(id64)
	if err != nil {
		log.Printf("%q\n", err)
		return
	}
	return
}

func (s *Sqlite) GetItems() (items []dbmodels.Item, err error) {

	items = make([]dbmodels.Item, 0)
	if err != nil {
		fmt.Println("error opening db")
		log.Fatal(err)
	}

	query := "SELECT id,url,item FROM Items;"
	rows, err := s.db.Query(query)
	defer rows.Close()
	for rows.Next() {
		item := dbmodels.Item{}
		err = rows.Scan(&item.Id, &item.Url, &item.Name)
		if err != nil {
			log.Println("error", err)
			continue
		} else {
			items = append(items, item)
		}
	}
	return
}

func (s *Sqlite) GetSettings() (settings dbmodels.Settings, err error) {

	query := `
			select refresh_interval,
			user_agent,
			discord_webhook,
			enabled from settings;
	`
	rows, err := s.db.Query(query)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&settings.Refresh_interval, &settings.User_agent, &settings.Discord_webhook, &settings.Enabled)
		if err != nil {
			log.Fatal(err)
		}

	}
	return
}

func (s *Sqlite) SaveSettings(settings dbmodels.Settings) (err error) {

	updatestmt := `
			update Settings set refresh_interval = ?,
			user_agent = ?,
			discord_webhook = ?,
			enabled = ?;
	`
	stmt, err := s.db.Prepare(updatestmt)
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
	return err
}
