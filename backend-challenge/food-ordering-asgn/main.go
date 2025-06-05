package main

import (
	"net/http"
	"os"

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

	requestHandler := pkg.NewRequestHandler(db)

	logger.Info("Starting server on port: 8080")
	if err := http.ListenAndServe(":8080", requestHandler.ServeHTTP()); err != nil {
		logger.Info("Failed to terminate server gracefully:", err)
		os.Exit(1)
	}
}