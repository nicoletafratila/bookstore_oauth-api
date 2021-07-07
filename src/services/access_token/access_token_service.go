package access_token

import (
	"github.com/nicoletafratila/bookstore_oauth-api/src/domain/access_token"
	"github.com/nicoletafratila/bookstore_oauth-api/src/repository/db"
	"github.com/nicoletafratila/bookstore_oauth-api/src/repository/rest"
	"github.com/nicoletafratila/bookstore_oauth-api/src/utils/errors"
	"strings"
)

//type Repository interface {
//	GetById(string) (*AccessToken, *errors.RestErr)
//	Create(AccessToken) *errors.RestErr
//	UpdateExpirationTime(AccessToken) *errors.RestErr
//}

type Service interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr)
	//CreateAccessToken(access_token.AccessToken)  *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type service struct {
	restRepository rest.RestUsersRepository
	dbRepository   db.DbRepository
}

func NewService(restRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restRepository: restRepo,
		dbRepository:   dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*access_token.AccessToken, *errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, errors.NewBadRequestError("invalid access token id")
	}
	accessToken, err := s.dbRepository.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

//func (s *service) CreateAccessToken(at access_token.AccessToken) *errors.RestErr {
//	if err := at.Validate(); err != nil {
//		return err
//	}
//	return s.dbRepository.Create(at)
//}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	//TODO: support both  grant types: client_credentials and password

	//Authenticate the user against the Users API
	user, err := s.restRepository.Login(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	//Generate a new access token
	at := access_token.GetNewAccessToken(user.Id)
	at.Generate()

	//Save the new access token in Cassandra
	if err := s.dbRepository.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepository.UpdateExpirationTime(at)
}
