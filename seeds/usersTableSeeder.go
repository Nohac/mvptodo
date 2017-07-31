package seeds

import (
	"fmt"

	"github.com/gocraft/dbr"
	"github.com/icrowley/fake"
	"github.com/noh4ck/mvptodo/models"
)

type usersTableSeeder struct{}

func (s usersTableSeeder) Seed(db *dbr.Session) {
	tx, _ := db.Begin()
	tx.Query("SET FOREIGN_KEY_CHECK = 0")
	tx.Query("TRUNCATE users;")

	fmt.Println("Run usersTableSeeder")

	for i := 0; i < 20; i++ {
		password, _ := models.PasswordHash("dritt")
		user := models.User{
			Email:    fake.EmailAddress(),
			Name:     fake.FullName(),
			Password: password,
		}

		tx.InsertInto("users").
			Columns("name", "email", "password").
			Record(user).Exec()
	}

	tx.Query("SET FOREIGN_KEY_CHECK = 1")

	err := tx.Commit()
	if err != nil {
		fmt.Println(err)
	}
}
