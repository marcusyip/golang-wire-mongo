package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	huhttp "github.com/marcusyip/golang-wire-mongo/core/http"
	huaccess "github.com/marcusyip/golang-wire-mongo/domains/access"
	huauth "github.com/marcusyip/golang-wire-mongo/domains/auth"
	ents "github.com/marcusyip/golang-wire-mongo/entities"
	"github.com/marcusyip/golang-wire-mongo/entities/errors"
	"github.com/marcusyip/golang-wire-mongo/models"
	mids "github.com/marcusyip/golang-wire-mongo/web/middlewares"
	"github.com/sirupsen/logrus"
)

type AccessController struct {
	logger    *logrus.Logger
	authMgr   *huauth.Manager
	accessMgr *huaccess.Manager
	accessEnt *ents.AccessEntity
}

func (c *AccessController) Create(ctx *gin.Context) {
	log := c.getLogger("Create")
	sessCtx := mids.SessionContext(ctx)

	params := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := ctx.BindJSON(params)
	if err != nil {
		log.WithField("err", err).Info(huhttp.LogInvalidJSON)
		ctx.Abort()
		return
	}
	authorizeParams := &huauth.AuthorizeParams{
		Username: params.Username,
		Password: params.Password,
	}
	user, err := c.authMgr.Authorize(sessCtx, authorizeParams)
	if err != nil {
		if err == huauth.ErrInvalidCredential {
			log.WithField("err", err).Info("Invalid credential")
			ctx.JSON(http.StatusUnauthorized, errors.InvalidCredential)
			return
		}
		log.WithField("err", err).Error("Failed to authorize with password")
		ctx.JSON(http.StatusInternalServerError, errors.InternalServerError)
		return
	}
	access, err := c.accessMgr.CreateForUser(sessCtx, user, models.AccessProviderFacebook)
	if err != nil {
		log.WithField("err", err).Error("Failed to create access for user")
		ctx.JSON(http.StatusInternalServerError, errors.InternalServerError)
		return
	}
	log.Info("Successfully create new access")
	ctx.JSON(http.StatusCreated, c.accessEnt.New(access))
}

func (c *AccessController) getLogger(method string) *logrus.Entry {
	return c.logger.WithFields(logrus.Fields{"controller": "api/access", "method": method})
}

func NewAccessController(
	logger *logrus.Logger,
	authMgr *huauth.Manager,
	accessMgr *huaccess.Manager,
	accessEnt *ents.AccessEntity,
) *AccessController {
	return &AccessController{
		logger:    logger,
		authMgr:   authMgr,
		accessMgr: accessMgr,
		accessEnt: accessEnt,
	}
}
