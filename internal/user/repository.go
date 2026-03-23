package user

import (
	"Avito/pkg/db"

	"github.com/google/uuid"
)

type UserRepository struct {
	DataBase *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		DataBase: database,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.DataBase.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) GetByID(id uuid.UUID) (*User, error) {
	var user User
	result := repo.DataBase.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := repo.DataBase.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}