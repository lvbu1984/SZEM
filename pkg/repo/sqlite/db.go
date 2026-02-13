package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// 基本连接检查
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

