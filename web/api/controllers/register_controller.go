package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	huhttp "github.com/marcusyip/golang-wire-mongo/core/http"
	huacct "github.com/marcusyip/golang-wire-mongo/domains/account"
	huauth "github.com/marcusyip/golang-wire-mongo/domains/auth"
	ents "github.com/marcusyip/golang-wire-mongo/entities"
	"github.com/marcusyip/golang-wire-mongo/entities/errors"
	mids "github.com/marcusyip/golang-wire-mongo/web/middlewares"
	"github.com/sirupsen/logrus"
)

type RegisterController struct {
	logger  *logrus.Logger
	authMgr *huauth.Manager
	acctMgr *huacct.Manager
	acctEnt *ents.AccountEntity
}

func (c *RegisterController) Register(ctx *gin.Context) {
	log := c.getLogger("Register")
	sessCtx := mids.SessionContext(ctx)
	params := &struct {
		Password    string `json:"password"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	}{}
	err := ctx.BindJSON(params)
	if err != nil {
		log.WithField("err", err).Info(huhttp.LogInvalidJSON)
		ctx.JSON(403, nil)
		return
	}
	validateParams := &huacct.ValidateParams{
		Username: &params.Username,
		Email:    &params.Email,
	}
	_, err = c.acctMgr.Validate(sessCtx, validateParams)
	if err != nil {
		if fieldErrors, ok := err.(huacct.FieldErrors); ok {
			fmt.Printf("############## %+v\n", fieldErrors)
			fieldError := fieldErrors.Errors()[0].(huacct.FieldError)
			switch fieldError.Err {
			case huacct.ErrUsernameAlreadyExists:
				ctx.JSON(http.StatusBadRequest, errors.UsernameAlreadyExists)
				return
			case huacct.ErrEmailAlreadyExists:
				ctx.JSON(http.StatusBadRequest, errors.EmailAlreadyExists)
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errors.InternalServerError)
		return
	}
	hashedPassword, err := c.authMgr.HashPassword(params.Password)
	if err != nil {
		return
	}
	createParams := &huacct.CreateParams{
		HashedPassword: &hashedPassword,
		Username:       &params.Username,
		Email:          params.Email,
		DisplayName:    params.DisplayName,
	}
	user, err := c.acctMgr.Create(sessCtx, createParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.InternalServerError)
		return
	}
	ctx.JSON(http.StatusCreated, c.acctEnt.New(user))
}

func (c *RegisterController) getLogger(method string) *logrus.Entry {
	return c.logger.WithFields(logrus.Fields{"controller": "api/register", "method": method})
}

func NewRegisterController(
	logger *logrus.Logger,
	authMgr *huauth.Manager,
	acctMgr *huacct.Manager,
	acctEnt *ents.AccountEntity,
) *RegisterController {
	return &RegisterController{
		logger:  logger,
		authMgr: authMgr,
		acctMgr: acctMgr,
		acctEnt: acctEnt,
	}
}
