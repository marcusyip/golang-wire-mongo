package controllers

import (
	"github.com/gin-gonic/gin"
	huhttp "github.com/marcusyip/golang-wire-mongo/core/http"
	huaccess "github.com/marcusyip/golang-wire-mongo/domains/access"
	huoauth "github.com/marcusyip/golang-wire-mongo/domains/oauth2"
	ents "github.com/marcusyip/golang-wire-mongo/entities"
	"github.com/marcusyip/golang-wire-mongo/models"
	mids "github.com/marcusyip/golang-wire-mongo/web/middlewares"
	"github.com/sirupsen/logrus"
)

type OauthFacebookController struct {
	logger    *logrus.Logger
	oauthMgr  *huoauth.Manager
	accessMgr *huaccess.Manager
	accessEnt *ents.AccessEntity
}

func (c *OauthFacebookController) Callback(ctx *gin.Context) {
	log := c.getLogger("Facebook")
	params := &huoauth.AuthorizeParams{}
	err := ctx.BindJSON(params)
	if err != nil {
		log.WithField("err", err).Info(huhttp.LogInvalidJSON)
		ctx.Abort()
		return
	}
	sessCtx := mids.SessionContext(ctx)
	user, err := c.oauthMgr.Authorize(sessCtx, "facebook", params)
	if err != nil {
		log.WithField("err", err).Info("Failed to authorize by facebook")
		return
	}
	access, err := c.accessMgr.CreateForUser(sessCtx, user, models.AccessProviderFacebook)
	if err != nil {
		log.WithField("err", err).Error("Failed to create access for user")
		return
	}
	log.Info("Successfully create new access")
	ctx.JSON(201, c.accessEnt.New(access))
}

func (c *OauthFacebookController) getLogger(method string) *logrus.Entry {
	return c.logger.WithFields(logrus.Fields{"controller": "api/oauth_facebook", "method": method})
}

func NewOauthFacebookController(
	logger *logrus.Logger,
	oauthMgr *huoauth.Manager,
	accessMgr *huaccess.Manager,
	accessEnt *ents.AccessEntity,
) *OauthFacebookController {
	return &OauthFacebookController{
		logger:    logger,
		oauthMgr:  oauthMgr,
		accessMgr: accessMgr,
		accessEnt: accessEnt,
	}
}
