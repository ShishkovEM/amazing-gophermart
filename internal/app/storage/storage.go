package storage

import (
	"fmt"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/exceptions"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/repository"
)

type Storage struct {
	Repo repository.CommonRepository
}

func NewStorage(database string, migrationsDir string) (*Storage, error) {
	if len(database) > 0 {
		NewStorage := repository.NewPostgresDB(database)
		err := NewStorage.MigrateToTheLatestSchema(database, migrationsDir)
		if err != nil {
			return nil, err
		}
		fmt.Println("Using PostgreSQL Database")
		return &Storage{
			Repo: NewStorage,
		}, nil
	}

	return &Storage{}, exceptions.ErrNoDatabaseDSN

}
