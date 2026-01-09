package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockMetricsService struct {
	calledAuthor string
	calledCtx    context.Context
	shouldError  bool
	response     *services.MetricsResponse
}

func (m *mockMetricsService) GetMetrics(ctx context.Context, author string) (*services.MetricsResponse, error) {
	m.calledAuthor = author
	m.calledCtx = ctx
	if m.shouldError {
		return nil, errors.New("service error")
	}
	return m.response, nil
}

func TestGetMetrics_Status200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBooks := []models.Book{
		{ID: 1, Name: "The Go Programming Language", Author: "Alan Donovan", UnitsSold: 5000, Price: 40},
		{ID: 2, Name: "Clean Code", Author: "Robert C. Martin", UnitsSold: 15000, Price: 50},
		{ID: 3, Name: "The Pragmatic Programmer", Author: "Andrew Hunt", UnitsSold: 13000, Price: 45},
	}

	mockService := &mockMetricsService{
		response: &services.MetricsResponse{
			Books:                mockBooks,
			MeanUnitsSold:        11000,
			CheapestBook:         "The Go Programming Language",
			BooksWrittenByAuthor: 1,
		},
	}
	handler := NewGetMetrics(mockService)

	r := gin.Default()
	r.GET("/", handler.Handle())

	req := httptest.NewRequest(http.MethodGet, "/?author=Alan+Donovan", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var resBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.NoError(t, err)

	assert.Contains(t, resBody, "books")
	assert.Contains(t, resBody, "mean_units_sold")
	assert.Contains(t, resBody, "cheapest_book")
	assert.Contains(t, resBody, "books_written_by_author")

	assert.Equal(t, 11000, int(resBody["mean_units_sold"].(float64)))
	assert.Equal(t, "The Go Programming Language", resBody["cheapest_book"])
	assert.Equal(t, 1, int(resBody["books_written_by_author"].(float64)))
}

func TestGetMetrics_AuthorParameterCorrectlyPassed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockMetricsService{
		response: &services.MetricsResponse{
			Books:                []models.Book{},
			MeanUnitsSold:        0,
			CheapestBook:         "",
			BooksWrittenByAuthor: 0,
		},
	}
	handler := NewGetMetrics(mockService)

	r := gin.Default()
	r.GET("/", handler.Handle())

	testCases := []struct {
		name           string
		queryParam     string
		expectedAuthor string
	}{
		{
			name:           "Author with space",
			queryParam:     "/?author=J.R.R.+Tolkien",
			expectedAuthor: "J.R.R. Tolkien",
		},
		{
			name:           "Author without space",
			queryParam:     "/?author=C.S.+Lewis",
			expectedAuthor: "C.S. Lewis",
		},
		// Nota: Este caso ya no es válido porque author es requerido
		// Se removió el test de "Empty author" ya que ahora retorna 400
		{
			name:           "Author with special characters",
			queryParam:     "/?author=Ursula+K.+Le+Guin",
			expectedAuthor: "Ursula K. Le Guin",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.queryParam, nil)
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, http.StatusOK, res.Code)
			assert.Equal(t, tc.expectedAuthor, mockService.calledAuthor)
			assert.NotNil(t, mockService.calledCtx)
		})
	}
}

func TestGetMetrics_Status500_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockMetricsService{
		shouldError: true,
	}
	handler := NewGetMetrics(mockService)

	r := gin.Default()
	r.GET("/", handler.Handle())

	req := httptest.NewRequest(http.MethodGet, "/?author=Test", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	var resBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to get metrics", resBody["error"])
}

func TestGetMetrics_Status400_InvalidQueryParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockMetricsService{
		response: &services.MetricsResponse{
			Books:                []models.Book{},
			MeanUnitsSold:        0,
			CheapestBook:         "",
			BooksWrittenByAuthor: 0,
		},
	}
	handler := NewGetMetrics(mockService)

	r := gin.Default()
	r.GET("/", handler.Handle())

	testCases := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "No query params - should return 400",
			url:        "/",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Valid author param",
			url:        "/?author=Test",
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, tc.wantStatus, res.Code)
		})
	}
}

func TestGetMetrics_ContextPropagation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockMetricsService{
		response: &services.MetricsResponse{
			Books:                []models.Book{},
			MeanUnitsSold:        0,
			CheapestBook:         "",
			BooksWrittenByAuthor: 0,
		},
	}
	handler := NewGetMetrics(mockService)

	r := gin.Default()
	r.GET("/", handler.Handle())

	req := httptest.NewRequest(http.MethodGet, "/?author=Test", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.NotNil(t, mockService.calledCtx)
	assert.Equal(t, req.Context(), mockService.calledCtx)
}
