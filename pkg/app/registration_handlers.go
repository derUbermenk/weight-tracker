package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response *GenericResponse
		var credentials *Credentials

		err := c.ShouldBindJSON(credentials)

		if err != nil {
			response.Status = false
			response.Message = "Handler Error"

			log.Printf("handler error: %v, err")
			c.JSON(http.StatusBadRequest, response)
		}

		// start checking
		if valid := s.userService.ValidatePassword(credentials.Password); !valid {
			response.Status = false
			response.Message = "Weak Password"

			c.JSON(http.StatusOK, response)
			c.Abort()
			return
		}

		exists, err := s.userService.UserExists(credentials.Email)

		if exists {
			response.Status = false
			response.Message = "User Exists"

			c.JSON(http.StatusOK, response)
			c.Abort()
			return
		}

		hashedPassword, err := s.userService.HashPassword(credentials.Password)

		if err != nil {
			response.Status = false
			response.Message = "Internal server error"

			c.JSON(http.StatusInternalServerError, response)
			c.Abort()
			return
		}

		_, err = s.userService.CreateUser(credentials.Email, hashedPassword)

		if err != nil {
			response.Status = false
			response.Message = "User creation failed"

			c.JSON(http.StatusOK, response)
			c.Abort()
			return
		}

		c.Set("credentials", credentials)
		c.Next()
	}
}
