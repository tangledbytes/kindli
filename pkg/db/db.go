package db

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

var (
	db       *sql.DB
	preloads = []string{}
)

// Setup sets up the database
func Setup(path string) {
	if db != nil {
		return
	}

	var err error
	db, err = sql.Open("sqlite", path)
	utils.ExitIfNotNil(err)
	db.SetMaxOpenConns(1)

	for _, query := range preloads {
		_, err = db.Exec(query)
		if err != nil {
			err = fmt.Errorf("failed executing query: %s: %s", query, err)
			utils.ExitIfNotNil(err)
		}
	}
}

// Instance returns the database instance
func Instance() *sql.DB {
	if db == nil {
		panic("db not initialized")
	}

	return db
}

// RegisterPreload takes in a query and invokes them just after
// initializing the database.
//
// This function might be useful for creating tables, etc.
func RegisterPreload(query string) {
	preloads = append(preloads, query)
}
