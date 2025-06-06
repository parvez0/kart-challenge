package pkg

// db.go file implements the database connection and management for the food ordering system.
// It provides functions to initialize the database, create tables, and manage connections.

import (

	"github.com/parvez0/food-ordering-asgn/utils"

	_ "github.com/mattn/go-sqlite3"
	gormsqlite "gorm.io/driver/sqlite"
	gorm "gorm.io/gorm"
)

var (
	logger = utils.GetLogger()
)

func WithSqliteInMemoryDB() func() (*gorm.DB, error) {
	return func() (*gorm.DB, error) {
		logger.Debugf("Setting up Sqlite InMemory database")
		db, err := gorm.Open(gormsqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			return nil, utils.WrapError(err, "Failed to setup Sqlite InMemory database")
		}
		return db, nil
	}
}

func NewDB(dbEngine func() (*gorm.DB, error)) (*gorm.DB, error) {
	return dbEngine()
}