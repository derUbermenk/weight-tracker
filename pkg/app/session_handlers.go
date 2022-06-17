package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*
			extract passed json to credentials struct
			credentials is password and email
			validate the credentials using authservice function

			if credential is invalid
				return an invalid credential response

			if valid
				return a valid credential response
				return data with jwt accesstoken and refreshtoken
		*/

		var cred Credentials
		err := c.ShouldBindJSON(&cred)
		if err != nil {
			handlerError(c, err, "Request")
		}

		is_cred_valid, err := s.authService.ValidateCredentials(cred.Email, cred.Password)
		if err != nil {
			handlerError(c, err, "Service")
		}

		if !is_cred_valid {
			c.JSON(
				http.StatusUnauthorized,
				&GenericResponse{
					Status:  false,
					Message: "Invalid Credentials",
				},
			)
		}

		access_token, err := s.authService.GenerateAccessToken(cred)
		if err != nil {
			handlerError(c, err, "Service")
			return
		}

		refresh_token, err := s.authService.GenerateRefreshToken(cred)
		if err != nil {
			handlerError(c, err, "Service")
			return
		}

		// return the tokens
		c.JSON(
			http.StatusOK,
			&GenericResponse{
				Status:  true,
				Message: "Signed in successfully",
				Data: &AuthResponse{
					AccessToken:  access_token,
					RefreshToken: refresh_token,
				},
			},
		)
	}
}

func (s *Server) ValidateAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (s *Server) ValidateRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (s *Server) RefreshAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
