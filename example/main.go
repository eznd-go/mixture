package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/eznd-go/mixture"
	_ "github.com/eznd-go/mixture/example/migrations"
)

func main() {
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

	_ = mixture.ApplyVerbose(db, "prod")
}
