package controllers

import (
	"github.com/gin-gonic/gin"
	huacct "github.com/marcusyip/golang-wire-mongo/domains/account"
	ents "github.com/marcusyip/golang-wire-mongo/entities"
	mids "github.com/marcusyip/golang-wire-mongo/web/middlewares"
)

type AccountController struct {
	acctMgr *huacct.Manager
	acctEnt *ents.AccountEntity
}

func (c *AccountController) Show(ctx *gin.Context) {
	user := mids.User(ctx)
	ctx.JSON(200, c.acctEnt.New(user))
}

func (c *AccountController) Create(ctx *gin.Context) {
	sessCtx := mids.SessionContext(ctx)
	createParams := &huacct.CreateParams{
		Email: "long3981@gmail.com",
	}
	user, err := c.acctMgr.Create(sessCtx, createParams)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(400, gin.H{})
		return
	}
	ctx.JSON(201, c.acctEnt.New(user))
}

func (c *AccountController) Update(ctx *gin.Context) {

}

func NewAccountController(
	acctMgr *huacct.Manager,
	acctEnt *ents.AccountEntity,
) *AccountController {
	return &AccountController{
		acctMgr: acctMgr,
		acctEnt: acctEnt,
	}
}
