package seeds

import (
	"fmt"

	"github.com/gocraft/dbr"
	"github.com/icrowley/fake"
	"github.com/noh4ck/mvptodo/models"
)

type todosTableSeeder struct{}

func (s todosTableSeeder) Seed(db *dbr.Session) {
	tx, _ := db.Begin()
	tx.Query("SET FOREIGN_KEY_CHECK = 0")
	tx.Query("TRUNCATE todos;")

	fmt.Println("Run todosTableSeeder")

	for i := 0; i < 20; i++ {
		for j := 0; j < 10; j++ {
			todo := models.Todo{
				UserID:      int64(i),
				Title:       fake.Title(),
				Description: fake.Paragraph(),
				Status:      false,
			}

			tx.InsertInto("todos").
				Columns("user_id", "title", "description", "status").
				Record(todo).Exec()
		}
	}

	tx.Query("SET FOREIGN_KEY_CHECK = 1")

	err := tx.Commit()
	if err != nil {
		fmt.Println(err)
	}
}
