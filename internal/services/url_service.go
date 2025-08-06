package services

import (
	"errors"
	"log/slog"
	"url_shortener/internal/storage"
	"url_shortener/internal/storage/postgres"
)

type UrlService interface {
	SaveURL(urlToSave string, alias string) error
	GetURL(alias string) (string, error)
	DeleteURL(alias string) error
}

type urlService struct {
	urlStorage *postgres.Storage
	log        *slog.Logger
}

func NewURLService(storage postgres.Storage, logger *slog.Logger) UrlService {
	return &urlService{urlStorage: &storage, log: logger}
}

func (c *urlService) SaveURL(urlToSave string, alias string) error {
	const fn = "services.url_service.SaveURL"
	log := c.log.With(
		slog.String("fn", fn),
	)

	if err := c.urlStorage.SaveURL(urlToSave, alias); err != nil {
		if errors.Is(err, storage.ErrURLExist) {
			log.Error("data already exists", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			return ErrURLAlreadyExists
		}
		log.Error("server error during saving the URL", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return err
	}

	return nil
}

func (c *urlService) GetURL(alias string) (string, error) {
	const fn = "services.url_service.GetURL"
	log := c.log.With(
		slog.String("fn", fn),
	)

	url, err := c.urlStorage.GetURL(alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Error("url with provided alias was not found", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			return "", ErrURLNotFound
		}
		log.Error("error trying to get a url", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return "", err
	}

	return url, nil
}

func (c *urlService) DeleteURL(alias string) error {
	const fn = "services.url_service.DeleteURL"
	log := c.log.With(
		slog.String("fn", fn),
	)

	if err := c.urlStorage.DeleteURL(alias); err != nil {
		log.Error("error trying to delete an alias", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return err
	}

	return nil
}
