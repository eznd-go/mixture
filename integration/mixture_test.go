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

func TestMigrationSuite(t *testing.T) {
	suite.Run(t, &migrationTestSuite{})
}

func (s *migrationTestSuite) SetupSuite()    {}
func (s *migrationTestSuite) TearDownSuite() {}
func (s *migrationTestSuite) BeforeTest(suite, test string) {
	gormLog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			//LogLevel: logger.Info,
			Colorful: false,
		},
	)

	s.db, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormLog,
	})
}
func (s *migrationTestSuite) AfterTest(suite, test string) {
	db, _ := s.db.DB()
	if db != nil {
		_ = db.Close()
	}
}

func (s *migrationTestSuite) Test_CreateTable() {
	migrations := testdata.CreateTable()
	mx := mixture.New(s.db)
	for r := range migrations {
		mx.Add(mixture.ForAnyEnv, &migrations[r])
	}
	err := mx.Apply(mixture.ForProduction)
	s.Assert().NoError(err)

	var num int64
	err = s.db.Model(testdata.User20220101{}).Count(&num).Error
	s.Assert().NoError(err)
	s.Assert().Equal(int64(0), num)
}

func (s *migrationTestSuite) Test_DropTable() {
	migrations := append(testdata.CreateTable(), testdata.DropTable()...)
	mx := mixture.New(s.db)
	for r := range migrations {
		mx.Add(mixture.ForAnyEnv, &migrations[r])
	}
	err := mx.Apply(mixture.ForProduction)
	s.Assert().NoError(err)

	var num int64
	err = s.db.Model(testdata.User20220101{}).Count(&num).Error
	s.Assert().EqualError(err, "no such table: users")
}

func (s *migrationTestSuite) Test_CreateBatch() {
	migrations := append(testdata.CreateTable(), testdata.CreateBatch()...)
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

func (s *migrationTestSuite) Test_DeleteBatch() {
	mx := mixture.New(s.db)
	migrations1 := append(testdata.CreateTable(), testdata.CreateBatch()...)
	for r := range migrations1 {
		mx.Add(mixture.ForAnyEnv, &migrations1[r])
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
	s.Assert().Equal(2, users[1].ID)
	s.Assert().Equal(3, users[2].ID)

	migrations2 := testdata.DeleteBatch()
	for r := range migrations2 {
		mx.Add(mixture.ForAnyEnv, &migrations2[r])
	}
	err = mx.Apply(mixture.ForProduction)
	s.Assert().NoError(err)

	err = s.db.Model(testdata.User20220101{}).Count(&num).Error
	s.Assert().NoError(err)
	s.Assert().Equal(int64(0), num)
}

func (s *migrationTestSuite) Test_Update() {
	migrations := append(append(testdata.CreateTable(), testdata.CreateBatch()...), testdata.Update()...)
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

	var user testdata.User20220101
	err = s.db.Model(testdata.User20220101{}).Order("id asc").First(&user).Error
	s.Assert().NoError(err)
	s.Assert().Equal("QWERTY1", user.Name)
}
