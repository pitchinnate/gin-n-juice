package routes

import (
	"gin-n-juice/routes/auth"
	"gin-n-juice/routes/members"
	"gin-n-juice/routes/middleware"
	"gin-n-juice/routes/teams"
	"gin-n-juice/routes/users"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	authRouter := r.Group("/auth")
	{
		authRouter.POST("/login", auth.PostLogin)
		authRouter.POST("/register", auth.PostRegister)
		authRouter.POST("/forgot", auth.PostForgot)
		authRouter.POST("/reset", auth.PostReset)
		authRouter.GET("/verify", auth.GetVerify)
	}

	authenticatedRouter := r.Group("")
	authenticatedRouter.Use(middleware.JwtAuth())
	{
		authenticatedRouter.POST("/auth/resend", auth.PostResend)

		verifiedEmailRouter := authenticatedRouter.Group("")
		verifiedEmailRouter.Use(middleware.EmailVerified())
		{
			userRouter := verifiedEmailRouter.Group("/users")
			{
				userRouter.GET("", users.GetList)
				userRouter.POST("", users.Create)
				userRouter.GET("/:id", users.GetSingle)
				userRouter.PUT("/:id", users.Update)
				userRouter.DELETE("/:id", users.Delete)
			}

			teamRouter := verifiedEmailRouter.Group("/teams")
			{
				teamRouter.GET("", teams.GetList)
				teamRouter.POST("", teams.Create)
				teamRouter.GET("/:id", teams.GetSingle)
				teamRouter.PUT("/:id", teams.Update)
				teamRouter.DELETE("/:id", teams.Delete)

				memberRouter := teamRouter.Group("/:id/members")
				memberRouter.Use(middleware.Team())
				{
					memberRouter.GET("", members.GetList)
					memberRouter.POST("", members.Create)
					memberRouter.GET("/:mid", members.GetSingle)
					memberRouter.PUT("/:mid", members.Update)
					memberRouter.DELETE("/:mid", members.Delete)
				}
			}
		}
	}

	return r
}
