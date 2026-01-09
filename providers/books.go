package providers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"educabot.com/bookshop/models"
)

type BooksProvider interface {
	GetBooks(ctx context.Context) []models.Book
}

type HTTPBooksProvider struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPBooksProvider(baseURL string) *HTTPBooksProvider {
	return &HTTPBooksProvider{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *HTTPBooksProvider) GetBooks(ctx context.Context) []models.Book {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL, nil)

	if err != nil {
		log.Printf("error creating request: %v", err)
		return []models.Book{}
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Printf("error making HTTP request: %v", err)
		return []models.Book{}
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("unexpected status code: %d", resp.StatusCode)
		return []models.Book{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
		return []models.Book{}
	}

	var books []models.Book
	if err := json.Unmarshal(body, &books); err != nil {
		log.Printf("error unmarshaling JSON: %v", err)
		return []models.Book{}
	}

	return books
}
