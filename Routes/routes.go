package Routes

import (
	"jwt-todo/auth-server/Handlers"
	"jwt-todo/auth-server/Middlewares"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          60 * time.Second,
		Credentials:     true,
		ValidateHeaders: true,
	}))
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))

	r.Use(sessions.Sessions("sessionKey", store))

	v1 := r.Group("/auth")
	{
		v1.POST("/login", Handlers.Login)
		v1.POST("/create", Handlers.CreateAccount)
		v1.POST("/todo", Handlers.CreateTodo)
		v1.GET("/logout", Middlewares.TokenAuthMiddleware(), Handlers.Logout)
		v1.POST("/token/refresh", Handlers.Refresh)

	}
	return r
}
