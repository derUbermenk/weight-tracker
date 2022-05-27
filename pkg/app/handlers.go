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
		var userID int
		var response = struct {
			Status string
			Data   string
			UserID int
		}{
			Status: "failed",
			UserID: userID,
		}

		err := c.ShouldBindJSON(&newUser)

		if err != nil {
			response.Data = err.Error()

			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		userID, err = s.userService.New(newUser)

		if err != nil {
			response.Data = err.Error()

			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response.Status = "success"
		response.Data = "user created"
		response.UserID = userID

		c.JSON(http.StatusCreated, response)
	}
}

func (s *Server) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		userID, err := strconv.Atoi(c.Param("userId"))
		var response = struct {
			Status string
			Data   string
			UserID int
		}{
			Status: "failed",
			UserID: userID,
		}

		if err != nil {
			response.Data = err.Error()
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		userID, err = s.userService.Delete(userID)

		if err != nil {
			response.Data = err.Error()
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		response.Status = "success"
		response.Data = "user deleted"
		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var updateUser api.UpdateUserRequest
		var user api.User
		var response = struct {
			Status string
			Data   string
			User   api.User
		}{
			Status: "failed",
		}

		err := c.ShouldBindJSON(&updateUser)

		if err != nil {
			response.Data = err.Error()

			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		user, err = s.userService.Update(updateUser)

		if err != nil {
			response.Data = err.Error()

			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response.Status = "success"
		response.Data = "user updated"
		response.User = user

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
