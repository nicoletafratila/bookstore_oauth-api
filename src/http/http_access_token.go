package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nicoletafratila/bookstore_oauth-api/src/domain/access_token"
	accessTokenServices "github.com/nicoletafratila/bookstore_oauth-api/src/services/access_token"
	"github.com/nicoletafratila/bookstore_utils-go/rest_errors"
	"net/http"
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(*gin.Context)
	//CreateAccessToken(*gin.Context)
}

type accessTokenHandler struct {
	service accessTokenServices.Service
}

func NewAccessTokenHandler(service accessTokenServices.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetById(c *gin.Context) {
	accessToken, err := handler.service.GetById(c.Param("access_token_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var request access_token.AccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	accessToken, err := handler.service.Create(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, accessToken)
}
