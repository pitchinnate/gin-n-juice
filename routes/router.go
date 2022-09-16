package routes

import (
	"gin-n-juice/routes/auth"
	"gin-n-juice/routes/middleware"
	"gin-n-juice/routes/users"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	authRouter := r.Group("/auth")
	authRouter.POST("/login", auth.PostLogin)
	authRouter.POST("/register", auth.PostRegister)
	authRouter.POST("/forgot", auth.PostForgot)
	authRouter.POST("/reset", auth.PostReset)

	authenticatedRouter := r.Group("")
	authenticatedRouter.Use(middleware.JwtAuth())

	userRouter := authenticatedRouter.Group("/users")
	userRouter.GET("", users.GetList)
	userRouter.POST("", users.Create)
	userRouter.GET("/:id", users.GetSingle)
	userRouter.PUT("/:id", users.Update)
	userRouter.DELETE("/:id", users.Delete)

	return r
}
