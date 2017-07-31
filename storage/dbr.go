package storage

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
)

func New(user, password, database string) *dbr.Connection {
	// db, err := sql.Open("mysql", "user:password@/dbname")
	// err = db.Ping()

	con, err := dbr.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s",
			user,
			password,
			database,
		), nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	return con
}
