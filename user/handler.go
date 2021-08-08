package user

import (
	"errors"
	"io/ioutil"
	"net/http"
)

type UserHandler struct {
	Controller UserController
}

var ErrUnprocessableInput = errors.New("unprocessable input")
var ErrUnexpectedError = errors.New("unexpected error")
var ErrNotFound = errors.New("not found")

func NewHandler(c UserController) {
	uh := UserHandler{Controller: c}
	http.Handle("/users", http.HandlerFunc(uh.handleUsers))
}

func (u UserHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	var (
		resp []byte
		err  error
	)
	switch r.Method {
	case http.MethodGet:
		email := r.URL.Query().Get("email")
		token := r.URL.Query().Get("token")
		resp, err = u.Controller.ValidateUserToken(email, token)
	case http.MethodPost:
		b, _ := ioutil.ReadAll(r.Body)
		resp, err = u.Controller.CreateUser(b)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
	if err != nil {
		handleErr(w, err)
	} else {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func handleErr(w http.ResponseWriter, err error) {
	switch err {
	case ErrUnexpectedError:
		w.WriteHeader(http.StatusInternalServerError)
	case ErrUnprocessableInput:
		w.WriteHeader(http.StatusBadRequest)
	case ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(err.Error()))
}
