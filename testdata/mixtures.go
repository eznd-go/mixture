package testdata

import (
	"github.com/ezn-go/mixture"
)

type User20220101 struct {
	ID       int
	Name     string `gorm:"unique;not null"`
	Email    string
	IsActive bool
}

func (s User20220101) TableName() string {
	return "users"
}

var users20220101 = []User20220101{
	{ID: 1, Name: "John Doe", Email: "john@doe.com", IsActive: true},
	{ID: 2, Name: "John Smith", Email: "john@smith.com", IsActive: true},
	{ID: 3, Name: "Blocked User", Email: "some@boo.com", IsActive: false},
}

type User20220102 struct {
	ID       int
	Name     string `gorm:"unique;not null"`
	Email    string
	Phone    string
	IsActive bool
}

func (s User20220102) TableName() string {
	return "users"
}

func CreateTable() []mixture.M {
	return []mixture.M{
		{
			ID:       "20220101-001",
			Migrate:  mixture.CreateTableM(&User20220101{}),
			Rollback: mixture.DropTableR(&User20220101{}),
		},
	}
}

func CreateBatch() []mixture.M {
	return []mixture.M{
		{
			ID:       "20220101-002",
			Migrate:  mixture.CreateBatchM(users20220101),
			Rollback: mixture.DeleteBatchR(users20220101),
		},
	}
}

func DeleteBatch() []mixture.M {
	return []mixture.M{
		{
			ID:       "20220101-003",
			Migrate:  mixture.DeleteBatchM(users20220101),
			Rollback: mixture.CreateBatchR(users20220101),
		},
	}
}

func DropTable() []mixture.M {
	return []mixture.M{
		{
			ID:       "20220101-004",
			Migrate:  mixture.DropTableM(users20220101[0]),
			Rollback: mixture.CreateBatchR(users20220101[0]),
		},
	}
}

func Update() []mixture.M {
	return []mixture.M{
		{
			ID:       "20220101-005",
			Migrate:  mixture.UpdateM("users", "id = 1", "name", "QWERTY1"),
			Rollback: mixture.UpdateR("users", "id = 1", "name", "John Doe"),
		},
	}
}

func Delete() []mixture.M {
	return []mixture.M{
		{
			ID:       "20220101-006",
			Migrate:  mixture.DeleteM("users", "id = ?", 1),
			Rollback: mixture.CreateBatchR([]User20220101{users20220101[0]}),
		},
	}
}
