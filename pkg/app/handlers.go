package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

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

// handles an error coming from a function in the api, which then
// is considered an internal server error
// or from a function within the handler which then maps to a Bad Request
// this then adds a json response for the given gin context
func handlerError(c *gin.Context, err error, source string) {
	log.Printf("%v", err)

	if source == "Request" {
		c.JSON(
			http.StatusBadRequest,
			&GenericResponse{
				Status:  false,
				Message: "Bad Request",
			},
		)

	} else if source == "Server" {
		c.JSON(
			http.StatusInternalServerError,
			&GenericResponse{
				Status:  false,
				Message: "Internal Server Error",
			},
		)
	}

	return
}
