package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/parvez0/food-ordering-asgn/utils"
)

// seeder.go provides functions for setting up DB and populate data for the first run.
// The current requiremnt is to fill order and product details. Read and fillup coupun codes

// SeedDatabase initializes the database with initial data
func SeedDatabase(db *gorm.DB) error {
	// Migrate all tables in correct order
	if err := db.AutoMigrate(&CouponSource{}, &Coupon{}); err != nil {
		return utils.WrapError(err, "failed to migrate CouponSource table")
	}
	if err := db.AutoMigrate(&Product{}, &Order{}, &OrderItem{}); err != nil {
		return utils.WrapError(err, "failed to migrate Product table")
	}

	// First seed products
	products, err := seedProductData(db)
	if err != nil {
		return utils.WrapError(err, "failed to seed product data")
	}

	// Then seed orders using the created products
	if len(products) == 0 {
		return utils.WrapError(nil, "no products available for seeding orders")
	}

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("Could not get caller info")
	}
	// Seed coupons from files  
	return seedCoupons(filepath.Join(filepath.Dir(file), "../data"), db)
}

func seedProductData(db *gorm.DB) ([]Product, error) {
	// Create initial products
	products := []Product{
		{
			Name:     "Margherita Pizza",
			Price:    12.99,
			Category: "Pizza",
		},
		{
			Name:     "Pepperoni Pizza",
			Price:    14.99,
			Category: "Pizza",
		},
		{
			Name:     "Caesar Salad",
			Price:    8.99,
			Category: "Salad",
		},
		{
			Name:     "Garlic Bread",
			Price:    4.99,
			Category: "Sides",
		},
		{
			Name:     "Chocolate Cake",
			Price:    6.99,
			Category: "Dessert",
		},
		{
			Name:     "Chicken Waffle",
			Price:    1.00,
			Category: "Waffle",
		},
	}

	// Create products in database
	if err := db.Create(&products).Error; err != nil {
		return nil, utils.WrapError(err, "failed to create new products")
	}

	// Fetch all products to get their IDs
	var results []Product
	if err := db.Find(&results).Error; err != nil {
		return nil, utils.WrapError(err, "failed to fetch product data from db")
	}

	return results, nil
}

func seedCoupons(dirPath string, db *gorm.DB) error {
	files, err := getFilesInDirectory(dirPath)
	if err != nil {
		return err
	}

	// Regex for matching only string coupons
	stringPattern := regexp.MustCompile(`\w{8,10}`)

	for _, file := range files {
		// Create source record for the file
		source := &CouponSource{
			Source: file,
		}
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(source).Error; err != nil {
			return utils.WrapError(err, "failed to create couponSource record")
		}

		// Read and process the file
		scanner, err := fileScanner(file)
		if err != nil {
			return utils.WrapError(err, "failed to read file: "+file)
		}

		// Process each line in the file
		for scanner.Scan() {
			couponCode := scanner.Text()
			if couponCode == "" || !stringPattern.MatchString(couponCode) {
				continue
			}

			// Create or find coupon record
			var coupon Coupon
			result := db.Where("code = ?", couponCode).First(&coupon)
			if result.Error != nil {
				if result.Error == gorm.ErrRecordNotFound {
					// Create new coupon
					coupon = Coupon{
						Code: couponCode,
					}
					if err := db.Create(&coupon).Error; err != nil {
						return utils.WrapError(err, "failed to create coupon record")
					}
				} else {
					return utils.WrapError(result.Error, "failed to check existing coupon")
				}
			}

			if err := db.Model(&coupon).Association("SourceFile").Append(source); err != nil {
				return utils.WrapError(err, "failed to create association for coupon: "+couponCode)
			}
		}

		if err := scanner.Err(); err != nil {
			return utils.WrapError(err, "error reading file: "+file)
		}
	}


	var coupons []Coupon
	if err := db.Preload("SourceFile").Find(&coupons).Error; err != nil {
		return utils.WrapError(err, "Failed to fetch coupons")
	}
	for _, c := range coupons {
		logger.Info(c.Code, c.SourceFile)
	}

	return nil
}

// fileScanner creates a new scanner for reading a file line by line
func fileScanner(filePath string) (*bufio.Scanner, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return nil, utils.WrapError(err, "failed to open file")
	}
	return bufio.NewScanner(fd), nil
}

// getFilesInDirectory returns a list of all files in the specified directory at 1 level
func getFilesInDirectory(dirPath string) ([]string, error) {
	var files []string

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, utils.WrapError(err, fmt.Sprintf("directory %s does not exist", dirPath))
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return utils.WrapError(err, fmt.Sprintf("error accessing path %s", path))
		}

		if info.IsDir() && path != dirPath {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, utils.WrapError(err, "error walking through directory")
	}

	return files, nil
} 