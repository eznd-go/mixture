package mixtures

import (
	"github.com/eznd-go/mixture"
)

func init() {
	type User struct {
		Id        int64  `json:"id" gorm:"primaryKey,autoIncrement"`
		Username  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Password  string `json:"password"`
		Salt      string `json:"salt"`
	}

	mx := &mixture.M{
		ID:       "0001",
		Migrate:  mixture.CreateTableM(&User{}),
		Rollback: mixture.DropTableR(&User{}),
	}

	mixture.Add(mixture.ForAnyEnv, mx)
}
