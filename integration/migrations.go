package integration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func GetTestMigrations() []gormigrate.Migration {
	return []gormigrate.Migration{
		{
			ID: "20220307-001",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					ID       int
					Name     string `gorm:"unique;not null"`
					Email    string
					IsActive bool
				}
				err := tx.AutoMigrate(&User{})
				if err != nil {
					return err
				}

				users := []User{
					{ID: 1, Name: "John Doe", Email: "john@doe.com", IsActive: true},
					{ID: 2, Name: "John Smith", Email: "john@smith.com", IsActive: true},
					{ID: 3, Name: "Blocked User", Email: "some@boo.com", IsActive: false},
				}
				for _, user := range users {
					err = tx.Create(&user).Error
					if err != nil {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
	}
}
