package web

import (
	"github.com/gin-gonic/gin"
	"github.com/marcusyip/golang-wire-mongo/web/api"
)

type Server struct {
	engine    *gin.Engine
	apiRouter *api.Router
	// oauthRouter *oauth.Router
}

func (s *Server) Start() {
	s.apiRouter.With(s.engine)
	// s.oauthRouter.With(s.engine)
	s.engine.Run()
}

func NewServer(
	engine *gin.Engine,
	apiRouter *api.Router,
	// oauthRouter *oauth.Router,
) *Server {
	return &Server{
		engine:    engine,
		apiRouter: apiRouter,
		// oauthRouter: oauthRouter,
	}
}
