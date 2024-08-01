package drivers

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SqLiteDriver struct {
	db *sql.DB
}

func NewSqLiteDriver() *SqLiteDriver {
	db, _ := sql.Open("sqlite3", "storage/database.db")

	query := `CREATE TABLE IF NOT EXISTS configs(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				group_key VARCHAR(35) NOT NULL,
				project_key VARCHAR(35) NOT NULL,
				payload TEXT NOT NULL
	 		)`

	_, err := db.Exec(query)
	if err != nil {
		log.Panic("Failed to initiate sqlite database")
	}

	return &SqLiteDriver{db: db}
}

func (c *SqLiteDriver) Exists(group string, key string) bool {
	var count int

	query := fmt.Sprintf(`SELECT COUNT(*) FROM configs WHERE group_key='%s' AND project_key='%s'`, group, key)
	c.db.QueryRow(query).Scan(&count)

	return count > 0
}

func (c SqLiteDriver) Show(group string, key string) string {
	var payload string

	query := fmt.Sprintf(`SELECT payload FROM configs WHERE group_key='%s' AND project_key='%s'`, group, key)
	c.db.QueryRow(query).Scan(&payload)

	return payload
}

func (c *SqLiteDriver) Update(group string, key string, payload string) string {

	query := fmt.Sprintf(`INSERT INTO configs (group_key,project_key,payload) VALUES ('%s','%s','%s')`, group, key, payload)
	if c.Exists(group, key) {
		query = fmt.Sprintf(`UPDATE configs SET payload='%s' WHERE group_key='%s' AND project_key='%s'`, payload, group, key)
	}

	_, err := c.db.Exec(query)

	if err != nil {
		return err.Error()
	}

	return payload
}
