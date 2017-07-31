package seeds

import "github.com/gocraft/dbr"

type Seeder interface {
	Seed(db *dbr.Session)
}

func Run(db *dbr.Session) {
	seeds := []Seeder{
		usersTableSeeder{},
		todosTableSeeder{},
	}

	for _, seed := range seeds {
		seed.Seed(db)
	}
}
