package main

import (
	"github.com/parvez0/food-ordering-asgn/pkg"
	"github.com/parvez0/food-ordering-asgn/utils"
)

func main() {
	logger := utils.GetLogger()
	
	db, err := pkg.NewDB(pkg.WithSqliteInMemoryDB())
	if err != nil {
		logger.Fatalf("Failed to setup database: %v", err)
	}
	logger.Infof("Database setup complete")
	err = pkg.SeedDatabase(db)
	if err != nil {
		logger.Fatalf("Failed to seed database: %v", err)
	}

	logger.Infof("Database migration complete")
}