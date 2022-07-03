package integration_test

import (
	"github.com/ezn-go/mixture/testdata"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ezn-go/mixture"
)

type migrationTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func TestMigrationTestSuite(t *testing.T) {
	suite.Run(t, &migrationTestSuite{})
}

func (s *migrationTestSuite) SetupSuite()    {}
func (s *migrationTestSuite) TearDownSuite() {}
func (s *migrationTestSuite) SetupTest() {
	gormLog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			//LogLevel:      logger.Silent,
			LogLevel: logger.Info,
			Colorful: false,
		},
	)

	s.db, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormLog,
	})
}
func (s *migrationTestSuite) TearDownTest() {}

func (s *migrationTestSuite) Test_Mixture_HappyPath() {
	migrations := testdata.GetHappyPathTestMigrations()
	mx := mixture.New(s.db)
	for r := range migrations {
		mx.Add(mixture.ForAnyEnv, &migrations[r])
	}
	err := mx.Apply(mixture.ForProduction)
	s.Assert().NoError(err)

	var num int64
	err = s.db.Model(testdata.User20220101{}).Count(&num).Error
	s.Assert().NoError(err)
	s.Assert().Equal(int64(3), num)

	var users []testdata.User20220101
	err = s.db.Model(testdata.User20220101{}).Order("id asc").Find(&users).Error
	s.Assert().NoError(err)
	s.Assert().Equal(3, len(users))
	s.Assert().Equal(1, users[0].ID)
}
