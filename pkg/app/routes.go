package app

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine {
	router := s.router

	// group all routes under /v1/api
	v1 := router.Group("/v1/api")
	{
		v1.GET("/status", s.ApiStatus())

		registration := router.Group("registration/")
		{
			registration.POST("/signUp", s.RegisterUser())
		}

		session := router.Group("session/")
		{
			session.GET("/logIn", s.LogIn())
			session.GET("/refreshToken", s.ValidateRefreshToken(), s.RefreshAccessToken())
		}
		// prefix the user routes

		private := router.Group("private/")
		private.Use(s.ValidateAccessToken())
		{
			user := private.Group("/user")
			{
				user.GET("/:userId", s.GetUser())       // show
				user.GET("", s.GetUsers())              // index
				user.POST("", s.CreateUser())           // create
				user.DELETE("/:userId", s.DeleteUser()) // delete
				user.PUT("/:userId", s.UpdateUser())    // update
			}

			weight := private.Group("/weight")
			{
				weight.POST("", s.CreateWeightEntry())
			}
		}
	}

	return router
}
