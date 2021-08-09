package user

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/onelogin/onelogin-go-sdk/pkg/services/users"

	"github.com/dcaponi/pw_less/cache"
	"github.com/dcaponi/pw_less/email"
	"github.com/google/uuid"
)

type UserController struct {
	Repo    UserRepository
	Cache   cache.Cache
	Emailer email.Emailer
}

func NewController(r UserRepository, c cache.Cache, e email.Emailer) UserController {
	return UserController{Repo: r, Cache: c, Emailer: e}
}

func (c UserController) ValidateUserToken(email, token string) ([]byte, error) {
	t, err := c.Cache.Get(email)
	if err != nil {
		log.Printf("invalid email %s given\n", email)
		return []byte("invalid email given"), ErrUnprocessableInput
	}
	if t == "" {
		log.Printf("email not known %s given\n", email)
		return []byte("email not known"), ErrUnprocessableInput
	}
	if token != t {
		log.Printf("invalid token %s given\n", token)
		return []byte("invalid token given"), ErrUnprocessableInput
	}

	u, err := c.Repo.GetByEmail(email)
	if err != nil {
		log.Printf("unable to retrieve user with email %s\n%v\n", email, err)
		return []byte("unexpected error"), ErrUnexpectedError
	}
	if u == nil {
		err = c.Repo.Create(email)
		if err != nil {
			log.Printf("unable to create user \n%v\n", err)
			return []byte("unable to create user"), ErrUnexpectedError
		}
		u, err = c.Repo.GetByEmail(email)
		if err != nil {
			log.Printf("unable to retrieve user with email %s\n%v\n", email, err)
			return []byte("unexpected error"), ErrUnexpectedError
		}
	}
	body, err := json.Marshal(u)
	if err != nil {
		log.Printf("unable parse repository output %s\n%v\n", email, err)
		return []byte("unexpected error"), ErrUnexpectedError
	}
	c.Cache.Del(email)
	return body, nil
}

func (c UserController) CreateUser(a []byte) ([]byte, error) {
	var u users.User

	if err := json.Unmarshal(a, &u); err != nil {
		log.Printf("invalid payload %v given\n", u)
		return []byte("invalid payload given"), ErrUnprocessableInput
	}

	existingUser, err := c.Repo.GetByEmail(*u.Email)
	if err != nil {
		log.Println("sql error", err)
	}

	if existingUser != nil {
		log.Println("existing user found", existingUser)
		token, err := c.Cache.Get(*existingUser.Email)
		if err != nil {
			token = strings.Replace(uuid.New().String(), "-", "", -1)
			c.Cache.Set(*existingUser.Email, token)
		}

		err = c.Emailer.Send([]string{*existingUser.Email}, fmt.Sprintf("http://localhost:8000/users?email=%s&token=%s", *existingUser.Email, token))
		if err != nil {
			log.Printf("unable to send email to user %s\n%v\n", *u.Email, err)
		}

		body, err := json.Marshal(existingUser)
		if err != nil {
			log.Printf("unable parse repository output %v\n%v\n", u, err)
			return []byte("unexpected error"), ErrUnexpectedError
		}
		return body, nil
	}

	body, err := json.Marshal(u)
	if err != nil {
		log.Printf("unable parse repository output %v\n%v\n", u, err)
		return []byte("unexpected error"), ErrUnexpectedError
	}

	token := strings.Replace(uuid.New().String(), "-", "", -1)
	c.Cache.Set(*u.Email, token)

	err = c.Emailer.Send([]string{*u.Email}, fmt.Sprintf("http://localhost:8000/users?email=%s&token=%s", *u.Email, token))
	if err != nil {
		log.Printf("unable to send email to user %s\n%v\n", *u.Email, err)
	}

	return body, nil
}
