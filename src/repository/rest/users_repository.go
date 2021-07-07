package rest

import (
	"encoding/json"
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/nicoletafratila/bookstore_oauth-api/src/domain/users"
	"github.com/nicoletafratila/bookstore_oauth-api/src/utils/errors"
	"time"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "https://localhost",
		Timeout: 100 * time.Microsecond,
	}
)

func NewRepository() RestUsersRepository {
	return &usersRepository{}
}

type RestUsersRepository interface {
	Login(string, string) (*users.User, *errors.RestErr)
}

type usersRepository struct {
}

func (r *usersRepository) Login(email string, password string) (*users.User, *errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response := usersRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, errors.NewInternalServerError("invalid rest client response when trying to login user")
	}

	if response.StatusCode > 299 {
		var restErr errors.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to login user")
		}
		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshal users login response")
	}

	return &user, nil
}
