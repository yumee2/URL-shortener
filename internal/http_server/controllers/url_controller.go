package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"url_shortener/internal/services"

	"github.com/gin-gonic/gin"
)

type UrlContoller interface {
	SaveURL(ctx *gin.Context)
	GetURL(ctx *gin.Context)
	DeleteURL(ctx *gin.Context)
}

type urlContoller struct {
	urlService services.UrlService
	log        *slog.Logger
}

type Request struct {
	URLToSave string `json:"urlToSave"`
	Alias     string `json:"alias"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func NewURLController(urlService services.UrlService, logger *slog.Logger) *urlContoller {
	return &urlContoller{urlService: urlService, log: logger}
}

func (c *urlContoller) SaveURL(ctx *gin.Context) {
	const fn = "controllers.url_controller.SaveURL"

	log := c.log.With(
		slog.String("fn", fn),
	)

	var requestJson Request
	if err := ctx.BindJSON(&requestJson); err != nil {
		log.Error("failed to parse json body", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if requestJson.URLToSave == "" {
		log.Error("missing required field", slog.String("field", "urlToSave"))
		ctx.JSON(400, gin.H{"error": "urlToSave is required"})
		return
	}

	if err := c.urlService.SaveURL(requestJson.URLToSave, requestJson.Alias); err != nil {
		if errors.Is(err, services.ErrURLAlreadyExists) {
			log.Error("data already exists", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			ctx.JSON(409, gin.H{"error": err.Error()})
			return
		}
		log.Error("server error during saving the URL", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"status": "OK"})
}

func (c *urlContoller) GetURL(ctx *gin.Context) {
	const fn = "controllers.url_controller.GetURL"

	log := c.log.With(
		slog.String("fn", fn),
	)

	alias := ctx.Param("alias")
	if alias == "" {
		log.Error("alias parameter is empty")
		ctx.JSON(400, gin.H{"error": "alias is required"})
		return
	}

	originalURL, err := c.urlService.GetURL(alias)
	if err != nil {
		if errors.Is(err, services.ErrURLNotFound) {
			log.Error("URL not found", slog.String("alias", alias))
			ctx.JSON(404, gin.H{"error": "URL not found"})
			return
		}
		log.Error("failed to retrieve URL", slog.String("error", err.Error()))
		ctx.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	ctx.Redirect(http.StatusFound, originalURL)
}

func (c *urlContoller) DeleteURL(ctx *gin.Context) {
	const fn = "controllers.url_controller.DeleteURL"

	log := c.log.With(
		slog.String("fn", fn),
	)

	alias := ctx.Param("alias")
	if alias == "" {
		log.Error("alias parameter is empty")
		ctx.JSON(400, gin.H{"error": "alias is required"})
		return
	}

	if err := c.urlService.DeleteURL(alias); err != nil {
		log.Error("error trying to delete the alias", slog.String("alias", alias))
		ctx.JSON(400, gin.H{"error": "error during deletign the url"})
		return
	}

	ctx.JSON(200, gin.H{
		"status": "OK",
	})
}
