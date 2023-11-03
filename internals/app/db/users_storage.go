package db

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"library/internals/app/models"
)

type UsersStorage struct {
	databasePool *pgxpool.Pool
}

func NewUsersStorage(pool *pgxpool.Pool) *UsersStorage {
	storage := new(UsersStorage)
	storage.databasePool = pool
	return storage
}

func (storage *UsersStorage) GetUsersList(nameFilter string) []models.User {
	query := "SELECT ID, name, email, age FROM users"
	args := make([]interface{}, 0)
	if nameFilter != "" {
		query += " WHERE name LIKE $1"
		args = append(args, fmt.Sprintf("%%%s%%", nameFilter))
	}
	var result []models.User
	err := pgxscan.Select(context.Background(), storage.databasePool, &result, query, args...)
	if err != nil {
		logrus.Errorln(err)
	}
	return result
}

func (storage *UsersStorage) GetUserById(id int64) models.User {
	query := "SELECT ID, name, email, age FROM users WHERE ID = $1"
	var result []models.User
	err := pgxscan.Select(context.Background(), storage.databasePool, &result, query, id)
	if err != nil {
		logrus.Errorln(err)
		return models.User{}
	}
	if len(result) == 0 {
		return models.User{}
	}
	return result[0]
}

func (storage *UsersStorage) CreateUser(user models.User) error {
	query := "INSERT INTO users(name,email, age) VALUES ($1,$2,$3)"
	_, err := storage.databasePool.Exec(context.Background(), query, user.Name, user.Email, user.Age)
	if err != nil {
		logrus.Errorln(err)
	}
	return nil
}
