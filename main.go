package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"

	"github.com/dcaponi/pw_less/cache"
	"github.com/dcaponi/pw_less/database"
	"github.com/dcaponi/pw_less/email"
	"github.com/dcaponi/pw_less/user"
)

func main() {
	db, err := database.New(database.DBConfig{
		Flavor:   os.Getenv("DATABASE_FLAVOR"),
		Host:     os.Getenv("DATABASE_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Db:       os.Getenv("POSTGRES_DB"),
	})
	if err != nil {
		log.Fatalln("failed to establish database connection!", err)
	}
	defer db.Close()

	cache, err := cache.NewRedisCache(cache.RedisCacheConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if err != nil {
		log.Fatalln("failed to establish cache connection!", err)
	}

	gmailer := email.SimpleEmailer{
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     os.Getenv("EMAIL_PORT"),
	}

	oneloginClient, err := client.NewClient(&client.APIClientConfig{
		Timeout:      client.DefaultTimeout,
		ClientID:     os.Getenv("ONELOGIN_CLIENT_ID"),
		ClientSecret: os.Getenv("ONELOGIN_CLIENT_SECRET"),
		Region:       os.Getenv("ONELOGIN_CLIENT_REGION"),
	})
	if err != nil {
		log.Fatalln("failed to establish onelogin connection!", err)
	}

	user.NewHandler(user.NewController(user.NewRepo(*oneloginClient), cache, gmailer))

	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
	if err != nil {
		log.Fatalln("unable to start server!", err)
	}
}
