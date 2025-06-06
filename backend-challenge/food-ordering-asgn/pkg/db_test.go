package pkg

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	suite "github.com/stretchr/testify/suite"
	gorm "gorm.io/gorm"
)

type DBTestSuite struct {
	suite.Suite
	dbInstance *gorm.DB
}

func (dbSuite *DBTestSuite) SetupTest() {
	var err error
	dbSuite.dbInstance, err = NewDB(WithSqliteInMemoryDB())

	assert.NoError(dbSuite.T(), err)

	err = dbSuite.dbInstance.AutoMigrate(
		&Product{},
		&Order{},
		&OrderItem{},
		&Coupon{},
		&CouponSource{},
	)
	
	assert.NoError(dbSuite.T(), err)
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

func (dbSuite *DBTestSuite) TestDBProductSeed() {
	product, err := seedProductData(dbSuite.dbInstance)
	assert.NoError(dbSuite.T(), err)
	assert.Equal(dbSuite.T(), len(product), 6)
}

func (dbSuite *DBTestSuite) TestCouponSeeding() {
	// Test coupon seeding
	_, file, _, ok := runtime.Caller(0)
	assert.True(dbSuite.T(), ok, "Expected to find caller file name")
	err := seedCoupons(filepath.Join(filepath.Dir(file), "../data"), dbSuite.dbInstance)
	assert.NoError(dbSuite.T(), err)

	// Verify coupons were created
	var coupons []Coupon
	err = dbSuite.dbInstance.Preload("SourceFile").Find(&coupons).Error
	assert.NoError(dbSuite.T(), err)
	assert.Greater(dbSuite.T(), len(coupons), 0)

	// Verify coupon sources were created
	var sources []CouponSource
	err = dbSuite.dbInstance.Find(&sources).Error
	assert.NoError(dbSuite.T(), err)
	assert.Greater(dbSuite.T(), len(sources), 0)

	// Verify coupon-source associations
	for _, coupon := range coupons {
		assert.Greater(dbSuite.T(), len(coupon.SourceFile), 0, "Coupon should have at least one source file")
	}

	// Test coupon validation by checking source files
	for _, coupon := range coupons {
		assert.GreaterOrEqual(dbSuite.T(), len(coupon.SourceFile), 1, "Coupon should have at least one source files to be valid")
	}
}