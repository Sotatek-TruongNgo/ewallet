package engine

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func NewEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS", "PUT", "PATCH"}
	corsConfig.AllowHeaders = []string{
		"RecaptchaToken",
		"AccessToken",
		"Authorization",
		"Content-Type",
		"Upgrade",
		"Origin",
		"Connection",
		"Accept-Encoding",
		"Accept-Language",
		"Host",
		"Access-Control-Request-Method",
		"Access-Control-Request-Headers",
		"official-account-id",
		"x-xss-protection",
	}

	server := gin.New()
	server.Use(cors.New(corsConfig))
	server.Use(gzip.Gzip(gzip.DefaultCompression))

	return server
}
