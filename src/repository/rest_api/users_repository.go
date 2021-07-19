package rest_api

import (
	"encoding/json"
	"errors"
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/nicoletafratila/bookstore_oauth-api/src/domain/users"
	"github.com/nicoletafratila/bookstore_utils-go/rest_errors"
	"time"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "https://localhost:8080",
		Timeout: 100 * time.Microsecond,
	}
)

func NewRepository() RestUsersRepository {
	return &usersRepository{}
}

type RestUsersRepository interface {
	Login(string, string) (*users.User, rest_errors.RestErr)
}

type usersRepository struct {
}

func (r *usersRepository) Login(email string, password string) (*users.User, rest_errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response := usersRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, rest_errors.NewInternalServerError("invalid rest_api client response when trying to login user", errors.New("rest client error"))
	}

	if response.StatusCode > 299 {
		apiErr, err := rest_errors.NewRestErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to login user", err)
		}
		return nil, apiErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal users login response", errors.New("rest client error"))
	}

	return &user, nil
}
