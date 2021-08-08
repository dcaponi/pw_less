package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Flavor, Host, User, Password, Port, Db, SSLMode string
}

func New(c DBConfig) (*sql.DB, error) {
	var psqlconn string
	if c.SSLMode != "" {
		psqlconn = fmt.Sprintf("postgres://%v:%v@%v:%v/%v", c.User, c.Password, c.Host, c.Port, c.Db)
	} else {
		psqlconn = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Db)
	}

	db, err := sql.Open(c.Flavor, psqlconn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Printf("established connection to %s!\n", c.Db)
	return db, nil
}
