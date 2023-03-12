package storage

import (
	"errors"
	"fmt"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/repository"
)

type Storage struct {
	Repo repository.CommonRepository
}

func NewStorage(database string) (*Storage, error) {
	if len(database) > 0 {
		NewStorage := repository.NewPostgresDB(database)
		fmt.Println("Using PostgreSQL Database")
		return &Storage{
			Repo: NewStorage,
		}, nil
	}

	return &Storage{}, errors.New("no database config found")

}
