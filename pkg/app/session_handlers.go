package app

import (
	"log"
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
		// get token string header
		// handle missing token
		// check if tkn is valid
		// if not
		// return internal server error

		// if not token valid
		// handle invalid token
		// otherwise move to next handler

		token_string := c.GetHeader("AccessToken")

		if token_string == "" {
			log.Print("Missing token")
			c.JSON(
				http.StatusBadRequest,
				&GenericResponse{Status: false, Message: "Missing token"},
			)

			c.Abort()
			return
		}

		token_status, current_user, err := s.authService.ValidateAccessToken(token_string)

		if err != nil {
			log.Printf("Internal Server Error: %v", err)

			c.JSON(
				http.StatusInternalServerError,
				&GenericResponse{Status: false, Message: "Server Error"},
			)

			c.Abort()
			return
		}

		switch token_status {

		case api.ExpiredAccessToken:
			// do something
			c.JSON(
				http.StatusUnauthorized,
				&GenericResponse{Status: false, Message: "Expired access token"},
			)
			c.Abort()
			return

		case api.TamperedAccessToken:
			c.JSON(
				http.StatusUnauthorized,
				&GenericResponse{Status: false, Message: "Tampered access token"},
			)
			c.Abort()
			return

		case api.ValidAccessToken:
			// do something
			c.Set("current_user", current_user)
			c.Next()
		default:
			// print invalid token
			log.Printf("Unknown Token Status: %v", token_status)
			c.JSON(
				http.StatusInternalServerError,
				&GenericResponse{Status: false, Message: "Internal server error"},
			)
			c.Abort()
			return
		}
	}
}

func (s *Server) ValidateRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get refresh token
		// check if refresh token exists
		// check if refresh token is valid
		// handle error

		// handle validity
		// set credentials for context
		// set refresh_token for context
		// move to next
		var creds Credentials

		refresh_token := c.GetHeader("RefreshToken")

		if refresh_token == "" {
			log.Printf("Missing refresh token")
			c.JSON(
				http.StatusBadRequest,
				&GenericResponse{Status: false, Message: "Missing refresh token"},
			)

			c.Abort()
			return
		}

		tkn_valid, email, err := s.authService.ValidateRefreshToken(refresh_token)

		if err != nil {
			log.Printf("Internal server error: %v", err)

			c.JSON(
				http.StatusInternalServerError,
				&GenericResponse{Status: false, Message: "Internal server error"},
			)
			c.Abort()
			return
		}

		if !tkn_valid {
			c.JSON(
				http.StatusUnauthorized,
				&GenericResponse{Status: false, Message: "invalid refresh token"},
			)
		}

		creds.Email = email

		c.Set("creds", creds)
		c.Set("refresh_token", refresh_token)
	}
}

func (s *Server) RefreshAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get returns an interface
		// .(Type) is type casting; casts the interface into the given type
		creds, _ := c.Get("creds")
		refresh_token, _ := c.Get("refresh_token")

		creds = creds.(Credentials)
		refresh_token = refresh_token.(string)

		access_token, err := s.authService.GenerateAccessToken(creds)

		if err != nil {
			log.Printf("Internal server error: %v", err)

			c.JSON(
				http.StatusInternalServerError,
				&GenericResponse{
					Status:  false,
					Message: "Internal server error",
				},
			)

			return
		}

		c.JSON(
			http.StatusOK,
			&GenericResponse{
				Status:  true,
				Message: "Access token refreshed",
				Data:    &AuthResponse{AccessToken: access_token, RefreshToken: refresh_token},
			},
		)
	}
}
