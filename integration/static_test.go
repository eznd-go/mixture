package integration_test

import (
	"github.com/ezn-go/mixture/integration"
	"log"
	"os"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ezn-go/mixture"
	"github.com/stretchr/testify/suite"
)

type staticTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func TestStaticTestSuite(t *testing.T) {
	suite.Run(t, &staticTestSuite{})
}

func (s *staticTestSuite) SetupSuite()    {}
func (s *staticTestSuite) TearDownSuite() {}
func (s *staticTestSuite) SetupTest() {
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
func (s *staticTestSuite) TearDownTest() {}

func (s *staticTestSuite) Test_Mixture_HappyPath() {
	migrations := integration.GetHappyPathTestMigrations()
	for r := range migrations {
		mixture.Add(mixture.ForAnyEnv, &migrations[r])
	}
	err := mixture.Apply(s.db, "prod")
	s.Assert().NoError(err)

	var num int64
	err = s.db.Model(integration.User{}).Count(&num).Error
	s.Assert().NoError(err)
	s.Assert().Equal(int64(3), num)

	var users []integration.User
	err = s.db.Model(integration.User{}).Order("id asc").Find(&users).Error
	s.Assert().NoError(err)
	s.Assert().Equal(3, len(users))
	s.Assert().Equal(1, users[0].ID)
}
