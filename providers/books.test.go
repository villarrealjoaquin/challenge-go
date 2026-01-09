package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"educabot.com/bookshop/models"
	"github.com/stretchr/testify/assert"
)

func TestHTTPBooksProvider_GetBooks_Success(t *testing.T) {
	expectedBooks := []models.Book{
		{ID: 1, Name: "The Fellowship of the Ring", Author: "J.R.R. Tolkien", UnitsSold: 50000000, Price: 20},
		{ID: 2, Name: "The Two Towers", Author: "J.R.R. Tolkien", UnitsSold: 30000000, Price: 20},
		{ID: 3, Name: "The Return of the King", Author: "J.R.R. Tolkien", UnitsSold: 50000000, Price: 20},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedBooks)
	}))
	defer server.Close()

	provider := NewHTTPBooksProvider(server.URL)
	ctx := context.Background()
	books := provider.GetBooks(ctx)

	assert.Len(t, books, 3)
	assert.Equal(t, expectedBooks[0].Name, books[0].Name)
	assert.Equal(t, expectedBooks[0].Author, books[0].Author)
	assert.Equal(t, expectedBooks[0].UnitsSold, books[0].UnitsSold)
	assert.Equal(t, expectedBooks[0].Price, books[0].Price)
}

func TestHTTPBooksProvider_GetBooks_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	provider := NewHTTPBooksProvider(server.URL)
	ctx := context.Background()
	books := provider.GetBooks(ctx)

	assert.Empty(t, books)
}

func TestHTTPBooksProvider_GetBooks_Non200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	provider := NewHTTPBooksProvider(server.URL)
	ctx := context.Background()
	books := provider.GetBooks(ctx)

	assert.Empty(t, books)
}

func TestHTTPBooksProvider_GetBooks_InternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	provider := NewHTTPBooksProvider(server.URL)
	ctx := context.Background()
	books := provider.GetBooks(ctx)

	assert.Empty(t, books)
}

func TestHTTPBooksProvider_GetBooks_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	provider := NewHTTPBooksProvider(server.URL)
	ctx := context.Background()
	books := provider.GetBooks(ctx)

	assert.Empty(t, books)
}

func TestHTTPBooksProvider_GetBooks_ServerUnreachable(t *testing.T) {
	provider := NewHTTPBooksProvider("http://localhost:99999/nonexistent")
	ctx := context.Background()
	books := provider.GetBooks(ctx)

	assert.Empty(t, books)
}

func TestHTTPBooksProvider_GetBooks_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	provider := NewHTTPBooksProvider(server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	books := provider.GetBooks(ctx)
	assert.Empty(t, books)
}

func TestHTTPBooksProvider_GetBooks_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "name": "Book"`))
	}))
	defer server.Close()

	provider := NewHTTPBooksProvider(server.URL)
	ctx := context.Background()
	books := provider.GetBooks(ctx)

	assert.Empty(t, books)
}

func TestHTTPBooksProvider_GetBooks_WrongJSONStructure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"not": "a book array"}`))
	}))
	defer server.Close()

	provider := NewHTTPBooksProvider(server.URL)
	ctx := context.Background()
	books := provider.GetBooks(ctx)

	assert.Empty(t, books)
}

func TestHTTPBooksProvider_GetBooks_ValidJSONWithAllFields(t *testing.T) {
	expectedBooks := []models.Book{
		{
			ID:        1,
			Name:      "The Fellowship of the Ring",
			Author:    "J.R.R. Tolkien",
			UnitsSold: 50000000,
			Price:     20,
		},
		{
			ID:        2,
			Name:      "The Two Towers",
			Author:    "J.R.R. Tolkien",
			UnitsSold: 30000000,
			Price:     20,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedBooks)
	}))
	defer server.Close()

	provider := NewHTTPBooksProvider(server.URL)
	ctx := context.Background()
	books := provider.GetBooks(ctx)

	assert.Len(t, books, 2)
	for i, book := range books {
		assert.Equal(t, expectedBooks[i].ID, book.ID)
		assert.Equal(t, expectedBooks[i].Name, book.Name)
		assert.Equal(t, expectedBooks[i].Author, book.Author)
		assert.Equal(t, expectedBooks[i].UnitsSold, book.UnitsSold)
		assert.Equal(t, expectedBooks[i].Price, book.Price)
	}
}
