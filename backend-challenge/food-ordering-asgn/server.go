package main

import (
	"github.com/parvez0/food-ordering-asgn/pkg"
	"github.com/parvez0/food-ordering-asgn/utils"
)

func main() {
	logger := utils.GetLogger()
	
	db, err := pkg.NewDB(pkg.WithSqliteInMemoryDB())
	if utils.IsNotNil(err) {
		logger.Fatalf("Failed to setup database: %v", err)
	}
	logger.Infof("Database setup complete")

	db.AutoMigrate(&pkg.Product{})
	db.AutoMigrate(&pkg.Order{})

	tables, err := db.Migrator().GetTables()
	if utils.IsNotNil(err) {
		logger.Fatal("failed to get tables")
	}
	logger.Info(tables)

	logger.Infof("Database migration complete")

}