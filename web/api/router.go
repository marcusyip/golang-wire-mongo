package api

import (
	"github.com/gin-gonic/gin"
	ctrls "github.com/marcusyip/golang-wire-mongo/web/api/controllers"
	mids "github.com/marcusyip/golang-wire-mongo/web/middlewares"
)

type Router struct {
	mgoSessMid *mids.MongoSession
	authMid    *mids.TokenAuthenticator

	regCtrl     *ctrls.RegisterController
	accessCtrl  *ctrls.AccessController
	oauthFbCtrl *ctrls.OauthFacebookController
	acctCtrl    *ctrls.AccountController
}

func (r *Router) With(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")
	v1.Use(r.mgoSessMid.Handler())

	public := r.newPublic(v1)

	private := r.newPrivate(v1)

	public.GET("/oauth/facebook", r.oauthFbCtrl.Callback)

	public.POST("/register", r.regCtrl.Register)
	public.POST("/access", r.accessCtrl.Create)

	private.GET("/account", r.acctCtrl.Show)
}

func (r *Router) newPublic(v1 *gin.RouterGroup) *gin.RouterGroup {
	public := v1.Group("/")
	public.Use(r.authMid.Handler(&mids.TokenAuthenticatorOptions{
		TokenRequired: false,
	}))
	return public
}

func (r *Router) newPrivate(v1 *gin.RouterGroup) *gin.RouterGroup {
	private := v1.Group("/")
	private.Use(r.authMid.Handler(&mids.TokenAuthenticatorOptions{
		TokenRequired: true,
	}))
	return private

}

func NewRouter(
	mgoSessMid *mids.MongoSession,
	authMid *mids.TokenAuthenticator,
	regCtrl *ctrls.RegisterController,
	accessCtrl *ctrls.AccessController,
	acctCtrl *ctrls.AccountController,
	oauthFbCtrl *ctrls.OauthFacebookController,
) *Router {
	router := &Router{
		mgoSessMid:  mgoSessMid,
		authMid:     authMid,
		regCtrl:     regCtrl,
		accessCtrl:  accessCtrl,
		acctCtrl:    acctCtrl,
		oauthFbCtrl: oauthFbCtrl,
	}
	return router
}
