package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var sessionContextKey = "session_context"

func SessionContext(ctx *gin.Context) mongo.SessionContext {
	return ctx.MustGet(sessionContextKey).(mongo.SessionContext)
}

func setMongoSession(ctx *gin.Context, session mongo.SessionContext) {
	ctx.Set(sessionContextKey, session)
}

type MongoSession struct {
	db *mongo.Database
}

func (m *MongoSession) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := m.db.Client().StartSession()
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
			setMongoSession(ctx, sessCtx)
			ctx.Next()
			return nil
		})
	}
}

func NewMongoSession(db *mongo.Database) *MongoSession {
	return &MongoSession{
		db: db,
	}
}
