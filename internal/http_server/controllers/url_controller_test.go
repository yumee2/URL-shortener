package controllers

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"url_shortener/internal/services"
	"url_shortener/internal/services/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(controller UrlContoller) *gin.Engine {
	router := gin.Default()
	router.POST("/url", controller.SaveURL)
	router.GET("/url/:alias", controller.GetURL)
	router.DELETE("/url/:alias", controller.DeleteURL)
	return router
}

func TestSaveURL(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
		mockSetup      func(*mocks.UrlService)
	}{
		{
			name:           "successful save",
			requestBody:    `{"urlToSave": "https://example.com", "alias": "test"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"status":"OK"}`,
			mockSetup: func(m *mocks.UrlService) {
				m.On("SaveURL", "https://example.com", "test").Return(nil)
			},
		},
		{
			name:           "missing urlToSave",
			requestBody:    `{"alias": "test"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"urlToSave is required"}`,
			mockSetup:      func(m *mocks.UrlService) {},
		},
		{
			name:           "invalid json",
			requestBody:    `invalid json`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid character 'i' looking for beginning of value"}`,
			mockSetup:      func(m *mocks.UrlService) {},
		},
		{
			name:           "url already exists",
			requestBody:    `{"urlToSave": "https://example.com", "alias": "test"}`,
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"alias already exists"}`,
			mockSetup: func(m *mocks.UrlService) {
				m.On("SaveURL", "https://example.com", "test").Return(services.ErrURLAlreadyExists)
			},
		},
		{
			name:           "server error",
			requestBody:    `{"urlToSave": "https://example.com", "alias": "test"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
			mockSetup: func(m *mocks.UrlService) {
				m.On("SaveURL", "https://example.com", "test").Return(errors.New("internal server error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := new(mocks.UrlService)
			tt.mockSetup(mockService)

			// Create controller with mock service
			controller := NewURLController(mockService, slog.Default())

			// Setup router
			router := setupRouter(controller)

			// Create request
			req, _ := http.NewRequest("POST", "/url", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetURL(t *testing.T) {
	tests := []struct {
		name           string
		alias          string
		expectedStatus int
		expectedBody   string
		mockSetup      func(*mocks.UrlService)
	}{
		{
			name:           "successful get",
			alias:          "test",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"originalURL":"https://example.com"}`,
			mockSetup: func(m *mocks.UrlService) {
				m.On("GetURL", "test").Return("https://example.com", nil)
			},
		},
		{
			name:           "url not found",
			alias:          "notfound",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"URL not found"}`,
			mockSetup: func(m *mocks.UrlService) {
				m.On("GetURL", "notfound").Return("", services.ErrURLNotFound)
			},
		},
		{
			name:           "server error",
			alias:          "test",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
			mockSetup: func(m *mocks.UrlService) {
				m.On("GetURL", "test").Return("", errors.New("internal server error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := new(mocks.UrlService)
			tt.mockSetup(mockService)

			// Create controller with mock service
			controller := NewURLController(mockService, slog.Default())

			// Setup router
			router := setupRouter(controller)

			// Create request
			req, _ := http.NewRequest("GET", "/url/"+tt.alias, nil)

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteURL(t *testing.T) {
	tests := []struct {
		name           string
		alias          string
		expectedStatus int
		expectedBody   string
		mockSetup      func(*mocks.UrlService)
	}{
		{
			name:           "successful delete",
			alias:          "test",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"OK"}`,
			mockSetup: func(m *mocks.UrlService) {
				m.On("DeleteURL", "test").Return(nil)
			},
		},
		{
			name:           "delete error",
			alias:          "test",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"error during deletign the url"}`,
			mockSetup: func(m *mocks.UrlService) {
				m.On("DeleteURL", "test").Return(errors.New("error during deletign the url"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := new(mocks.UrlService)
			tt.mockSetup(mockService)

			// Create controller with mock service
			controller := NewURLController(mockService, slog.Default())

			// Setup router
			router := setupRouter(controller)

			// Create request
			req, _ := http.NewRequest("DELETE", "/url/"+tt.alias, nil)

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}
