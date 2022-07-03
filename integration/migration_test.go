package integration_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ezn-go/mixture"
	"github.com/ezn-go/mixture/integration"
)

type migrationTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func TestMixtureTestSuite(t *testing.T) {
	suite.Run(t, &migrationTestSuite{})
}

func (s *migrationTestSuite) SetupSuite()    {}
func (s *migrationTestSuite) TearDownSuite() {}
func (s *migrationTestSuite) SetupTest() {
	gormLog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      false,
		},
	)

	s.db, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormLog,
	})
}
func (s *migrationTestSuite) TearDownTest() {}

type User struct {
	ID       int
	Name     string `gorm:"unique;not null"`
	Email    string
	IsActive bool
}

func (s *migrationTestSuite) Test_Mixture_HappyPath() {
	migrations := integration.GetHappyPathTestMigrations()
	mx := mixture.New(s.db)
	for r := range migrations {
		mx.Add(mixture.ForAnyEnv, &migrations[r])
	}
	err := mx.Apply(mixture.ForProduction)
	s.Assert().NoError(err)

	var num int64
	err = s.db.Model(User{}).Count(&num).Error
	s.Assert().NoError(err)
	s.Assert().Equal(int64(3), num)

	var users []User
	err = s.db.Model(User{}).Order("id asc").Find(&users).Error
	s.Assert().NoError(err)
	s.Assert().Equal(3, len(users))
	s.Assert().Equal(1, users[0].ID)
}
