package drivers

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SqLiteDriver struct {
	db *sql.DB
}

func NewSqLiteDriver() *SqLiteDriver {
	if err := os.MkdirAll("storage", 0755); err != nil {
		log.Panicf("Failed to create sqlite storage directory: %v", err)
	}

	db, err := sql.Open("sqlite3", "storage/database.db")
	if err != nil {
		log.Panicf("Failed to open sqlite database: %v", err)
	}

	query := `CREATE TABLE IF NOT EXISTS configs(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				group_key VARCHAR(35) NOT NULL,
				project_key VARCHAR(35) NOT NULL,
				payload TEXT NOT NULL
	 		)`

	_, err = db.Exec(query)
	if err != nil {
		log.Panicf("Failed to initiate sqlite database: %v", err)
	}

	query = `CREATE TABLE IF NOT EXISTS locks(
				group_key VARCHAR(35) NOT NULL,
				project_key VARCHAR(35) NOT NULL,
				payload TEXT NOT NULL,
				PRIMARY KEY(group_key, project_key)
			)`

	_, err = db.Exec(query)
	if err != nil {
		log.Panicf("Failed to initiate sqlite database: %v", err)
	}

	return &SqLiteDriver{db: db}
}

func (c *SqLiteDriver) Exists(group string, key string) (bool, error) {
	var count int

	err := c.db.QueryRow(`SELECT COUNT(*) FROM configs WHERE group_key = ? AND project_key = ?`, group, key).Scan(&count)

	return count > 0, err
}

func (c SqLiteDriver) Show(group string, key string) ([]byte, error) {
	var payload string

	err := c.db.QueryRow(`SELECT payload FROM configs WHERE group_key = ? AND project_key = ?`, group, key).Scan(&payload)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	return []byte(payload), err
}

func (c *SqLiteDriver) Update(group string, key string, payload []byte) ([]byte, error) {
	result, err := c.db.Exec(
		`UPDATE configs SET payload = ? WHERE group_key = ? AND project_key = ?`,
		string(payload),
		group,
		key,
	)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		_, err = c.db.Exec(
			`INSERT INTO configs (group_key, project_key, payload) VALUES (?, ?, ?)`,
			group,
			key,
			string(payload),
		)
		if err != nil {
			return nil, err
		}
	}

	return payload, nil
}

func (c *SqLiteDriver) Lock(group string, key string, payload []byte) ([]byte, error) {
	_, err := c.db.Exec(`INSERT INTO locks (group_key, project_key, payload) VALUES (?, ?, ?)`, group, key, string(payload))
	if err == nil {
		return payload, nil
	}

	var current string
	readErr := c.db.QueryRow(`SELECT payload FROM locks WHERE group_key = ? AND project_key = ?`, group, key).Scan(&current)
	if readErr != nil {
		return nil, err
	}

	return []byte(current), ErrLocked
}

func (c *SqLiteDriver) Unlock(group string, key string, payload []byte) error {
	var current string
	err := c.db.QueryRow(`SELECT payload FROM locks WHERE group_key = ? AND project_key = ?`, group, key).Scan(&current)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	if !sameLockID([]byte(current), payload) {
		return ErrUnlockMismatch
	}

	_, err = c.db.Exec(`DELETE FROM locks WHERE group_key = ? AND project_key = ?`, group, key)
	return err
}
