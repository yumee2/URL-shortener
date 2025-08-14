package routers

import (
	"url_shortener/internal/config"
	"url_shortener/internal/http_server/controllers"

	"github.com/gin-gonic/gin"
)

func SetupURLRoutes(r *gin.Engine, urlController controllers.UrlContoller, cfg config.Config) {
	urlGroup := r.Group("/url")

	urlGroup.GET("/:alias", urlController.GetURL)
	secured := urlGroup.Group("/", gin.BasicAuth(gin.Accounts{
		cfg.HttpServer.User: cfg.HttpServer.Password,
	}))
	{
		secured.POST("/", urlController.SaveURL)
		secured.DELETE("/:alias", urlController.DeleteURL)
	}
}
