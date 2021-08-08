package user

import (
	"database/sql"
)

type UserRepository struct {
	Store *sql.DB
}

func NewRepo(s *sql.DB) UserRepository {
	return UserRepository{Store: s}
}

func (r UserRepository) List() ([]User, error) {
	results := []User{}
	stmt := `select * from "users"`
	rows, err := r.Store.Query(stmt)
	if err != nil {
		return results, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var email string
		err = rows.Scan(&id, &email)
		if err != nil {
			return results, err
		}
		results = append(results, User{ID: id, Email: email})
	}
	return results, nil
}

func (r UserRepository) GetById(id int64) (User, error) {
	result := User{}
	stmt := `select * from "users" where "id"=$1`
	rows, err := r.Store.Query(stmt, id)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var email string
		err = rows.Scan(&id, &email)
		if err != nil {
			return result, err
		}
		result = User{ID: id, Email: email}
	}
	return result, nil
}

func (r UserRepository) GetByEmail(email string) (User, error) {
	result := User{}
	stmt := `select * from "users" where "email"=$1`
	rows, err := r.Store.Query(stmt, email)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var email string
		err = rows.Scan(&id, &email)
		if err != nil {
			return result, err
		}
		result = User{ID: id, Email: email}
	}
	return result, nil
}

func (r UserRepository) Create(u *User) error {
	stmt := `insert into "users"("email") values($1)`
	result, err := r.Store.Exec(stmt, u.Email)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	u.ID = id
	return nil
}

func (r UserRepository) Delete(id int64) error {
	stmt := `delete from "users" where id=$1`
	_, e := r.Store.Exec(stmt, id)
	return e
}
