package controllers

import (
	"github.com/gin-gonic/gin"
	huhttp "github.com/marcusyip/golang-wire-mongo/core/http"
	huauth "github.com/marcusyip/golang-wire-mongo/domains/auth"
	ents "github.com/marcusyip/golang-wire-mongo/entities"
	mids "github.com/marcusyip/golang-wire-mongo/web/middlewares"
	"github.com/sirupsen/logrus"
)

type RegisterController struct {
	logger  *logrus.Logger
	authMgr *huauth.Manager
	acctEnt *ents.AccountEntity
}

func (c *RegisterController) Register(ctx *gin.Context) {
	log := c.getLogger("Register")
	sessCtx := mids.MongoSession(ctx)
	params := &struct {
		Password    string `json:"password"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	}{}
	err := ctx.BindJSON(params)
	if err != nil {
		log.WithField("err", err).Info(huhttp.LogInvalidJSON)
		return
	}
	regParams := &huauth.RegisterParams{
		Password:    params.Password,
		Username:    params.Username,
		Email:       params.Email,
		DisplayName: params.DisplayName,
	}
	user, err := c.authMgr.Register(sessCtx, regParams)
	if err != nil {
		return
	}
	ctx.JSON(201, c.acctEnt.New(user))
}

func (c *RegisterController) getLogger(method string) *logrus.Entry {
	return c.logger.WithFields(logrus.Fields{"controller": "api/register", "method": method})
}

func NewRegisterController(
	logger *logrus.Logger,
	authMgr *huauth.Manager,
	acctEnt *ents.AccountEntity,
) *RegisterController {
	return &RegisterController{
		logger:  logger,
		authMgr: authMgr,
		acctEnt: acctEnt,
	}
}
