package user

import (
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/users"
)

type UserRepository struct {
	Store client.APIClient
}

func NewRepo(c client.APIClient) UserRepository {
	return UserRepository{Store: c}
}

func (r UserRepository) List() ([]users.User, error) {
	return r.Store.Services.UsersV2.Query(nil)
}

func (r UserRepository) GetById(id int32) (*users.User, error) {
	return r.Store.Services.UsersV2.GetOne(id)
}

func (r UserRepository) GetByEmail(email string) (*users.User, error) {
	u, err := r.Store.Services.UsersV2.Query(&users.UserQuery{Email: &email})
	if err != nil {
		return nil, err
	}
	if len(u) == 0 {
		return nil, nil
	}
	return &u[0], nil
}

func (r UserRepository) Create(email string) error {
	return r.Store.Services.UsersV2.Create(&users.User{Email: &email})
}

func (r UserRepository) Delete(id int32) error {
	return r.Store.Services.UsersV2.Destroy(id)
}
