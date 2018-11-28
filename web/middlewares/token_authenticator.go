package middlewares

import (
	"github.com/gin-gonic/gin"
	huaccess "github.com/marcusyip/golang-wire-mongo/domains/access"
	ents "github.com/marcusyip/golang-wire-mongo/entities"
	"github.com/marcusyip/golang-wire-mongo/entities/errors"
	"github.com/marcusyip/golang-wire-mongo/models"
	repos "github.com/marcusyip/golang-wire-mongo/repositories"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

var accessKey = "access"

func Access(ctx *gin.Context) *models.Access {
	return ctx.MustGet(accessKey).(*models.Access)
}

func setAccess(ctx *gin.Context, access *models.Access) {
	ctx.Set(accessKey, access)
}

var userKey = "user"

func User(ctx *gin.Context) *models.User {
	return ctx.MustGet(userKey).(*models.User)
}

func setUser(ctx *gin.Context, user *models.User) {
	ctx.Set(userKey, user)
}

type TokenAuthenticatorOptions struct {
	TokenRequired bool
}

type TokenAuthenticator struct {
	logger    *logrus.Logger
	accessMgr *huaccess.Manager
	userRepo  *repos.UserRepository
	errEnt    *ents.ErrorEntity
}

func (a *TokenAuthenticator) Handler(options *TokenAuthenticatorOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := a.getLogger()
		sessCtx := SessionContext(ctx)
		token := ctx.GetHeader("X-Access-Token")
		if len(token) == 0 {
			if options.TokenRequired {
				ctx.AbortWithStatusJSON(401, errors.UnauthorizedError)
				log.Info("Missing access token")
				return
			}
			ctx.Next()
			return
		}
		access, err := a.accessMgr.Resolve(sessCtx, token)
		if err != nil && err != mongo.ErrNoDocuments {
			ctx.AbortWithStatusJSON(500, errors.InternalServerError)
			log.Error("Cannot resolve access token")
			return
		}
		if access == nil {
			log.Info("Access token not found")
			ctx.Next()
			return
		}
		setAccess(ctx, access)
		log = log.WithField("access_id", access.ID)
		user, err := a.userRepo.FindByID(sessCtx, access.UserID)
		if err != nil {
			ctx.JSON(500, errors.InternalServerError)
			log.Info("Cannot find access's user")
			return
		}
		setUser(ctx, user)
		ctx.Next()
	}
}

func (a *TokenAuthenticator) getLogger() *logrus.Entry {
	return a.logger.WithFields(logrus.Fields{"middleware": "token_authenticator"})
}

func NewTokenAuthenticator(
	logger *logrus.Logger,
	accessMgr *huaccess.Manager,
	userRepo *repos.UserRepository,
	errEnt *ents.ErrorEntity,
) *TokenAuthenticator {
	return &TokenAuthenticator{
		logger:    logger,
		accessMgr: accessMgr,
		userRepo:  userRepo,
		errEnt:    errEnt,
	}
}
