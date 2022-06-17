package app

import "github.com/gin-gonic/gin"

func (s *Server) LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
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
