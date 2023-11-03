package processors

import (
	"errors"
	"library/internals/app/db"
	"library/internals/app/models"
)

type UsersProcessor struct {
	storage *db.UsersStorage
}

func NewUsersProcessor(storage *db.UsersStorage) *UsersProcessor {
	processor := new(UsersProcessor)
	processor.storage = storage
	return processor
}

func (processor *UsersProcessor) CreateUser(user models.User) error {
	if user.Name == "" {
		return errors.New("name should not be empty")
	}
	if user.Email == "" {
		return errors.New("email should not be empty")

	}
	if user.Age <= 0 {
		return errors.New("age must be greater than 0")
	}
	return processor.storage.CreateUser(user)
}

func (processor *UsersProcessor) FindUser(id int64) (models.User, error) {
	user := processor.storage.GetUserById(id)
	if user.Id != id {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (processor *UsersProcessor) ListUsers(nameFilter string) ([]models.User, error) {
	return processor.storage.GetUsersList(nameFilter), nil
}
