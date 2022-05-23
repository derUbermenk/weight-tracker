package app

import (
	"log"
	"net/http"
	"strconv"
	"weight-tracker/pkg/api"

	"github.com/gin-gonic/gin"
)

func (s *Server) ApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]string{
			"status": "success",
			"data":   "weight tracker API running smoothly",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("userId"))

		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		user, err := s.userService.GetUser(userID)

		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func (s *Server) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := s.userService.All()

		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func (s *Server) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var newUser api.NewUserRequest

		err := c.ShouldBindJSON(&newUser)

		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		err = s.userService.New(newUser)

		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		response := map[string]string{
			"status": "success",
			"data":   "new user created",
		}

		c.JSON(http.StatusCreated, response)
	}
}

func (s *Server) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var updateUser api.UpdateUserRequest

		err := c.ShouldBindJSON(&updateUser)

		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		err = s.userService.Update(updateUser)

		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		response := map[string]string{
			"status": "success",
			"data":   "user updated",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) CreateWeightEntry() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var newWeight api.NewWeightRequest
		err := c.ShouldBindJSON(&newWeight)

		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		err = s.weightService.New(newWeight)

		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		response := map[string]string{
			"status": "success",
			"data":   "new user created",
		}

		c.JSON(http.StatusOK, response)
	}
}
