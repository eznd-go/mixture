package mixture

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
	"time"
)

type User0001 struct {
	Id        int64  `json:"id" gorm:"primaryKey,autoIncrement"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Salt      string `json:"salt"`
}

func (u User0001) TableName() string {
	return "users"
}

type User0002 struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Salt      string `json:"salt"`
}

func (u User0002) TableName() string {
	return "users"
}

func Test_Verbose(t *testing.T) {
	gormLog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      false,
		},
	)

	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormLog,
	})

	m := New(db)

	m1 := &M{
		ID:       "0001",
		Migrate:  CreateTableM(&User0001{}),
		Rollback: DropTableR(&User0001{}),
	}
	m.Add(ForAnyEnv, m1)

	users := []User0002{
		{Username: "user1", FirstName: "User", LastName: "One", Email: "1@user.com", Phone: "+1", Salt: "", Password: ""},
		{Username: "user2", FirstName: "User", LastName: "Two", Email: "2@user.com", Phone: "+2", Salt: "", Password: ""},
		{Username: "user3", FirstName: "User", LastName: "Three", Email: "3@user.com", Phone: "+3", Salt: "", Password: ""},
	}

	m2 := &M{
		ID:       "0002",
		Migrate:  CreateBatchM(users),
		Rollback: DeleteBatchR(users),
	}
	m.Add(ForAnyEnv, m2)

	_ = m.ApplyVerbose(ForAnyEnv)
}
