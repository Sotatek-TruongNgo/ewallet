package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	doc "github.com/truongnqse05461/ewallet/generate/doc"
)

type SwaggerHandler struct {
}

func NewSwaggerHandler(
	host string,
	basePath string,
) *SwaggerHandler {
	doc.SwaggerInfo.Host = host
	doc.SwaggerInfo.BasePath = basePath
	return &SwaggerHandler{}
}

// @Tags swagger
// @Summary swagger doc entry
// @Id swagger
// @Router /swagger/index.html [get]
// @version 1.0
// @Success 200
func (*SwaggerHandler) Swagger() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
