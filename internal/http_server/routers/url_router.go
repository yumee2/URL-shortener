package routers

import (
	"url_shortener/internal/http_server/controllers"

	"github.com/gin-gonic/gin"
)

func SetupURLRoutes(r *gin.Engine, urlController controllers.UrlContoller) {
	urlGroup := r.Group("/url")
	{
		urlGroup.POST("/", urlController.SaveURL)
		urlGroup.GET("/:alias", urlController.GetURL)
		urlGroup.DELETE("/:alias", urlController.DeleteURL)
	}
}
