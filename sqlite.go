package main

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

func initializeDatabase() *sql.DB {
	db, _ := sql.Open("sqlite3", "database.db")

	query := `CREATE TABLE IF NOT EXISTS events (
				id INTEGER PRIMARY KEY,  
				source TEXT NOT NULL,
				payload JSON NOT NULL
	 		)`

	db.Exec(query)

	return db
}

func loadConfigFromDB(c echo.Context) {
	// group := c.Request().PathValue("group")
	// key := c.Request().PathValue("key")

	db := initializeDatabase()

	db.Close()

}
