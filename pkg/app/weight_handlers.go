package app

import (
	"log"
	"net/http"
	"weight-tracker/pkg/api"

	"github.com/gin-gonic/gin"
)

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
