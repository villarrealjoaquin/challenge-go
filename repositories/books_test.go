package repositories

import (
	"context"
	"testing"

	"educabot.com/bookshop/models"
	"github.com/stretchr/testify/assert"
)

type mockBooksProvider struct {
	books []models.Book
	ctx   context.Context
}

func (m *mockBooksProvider) GetBooks(ctx context.Context) []models.Book {
	m.ctx = ctx
	return m.books
}

func TestBooksRepository_GetAll_Success(t *testing.T) {
	expectedBooks := []models.Book{
		{ID: 1, Name: "The Fellowship of the Ring", Author: "J.R.R. Tolkien", UnitsSold: 50000000, Price: 20},
		{ID: 2, Name: "The Two Towers", Author: "J.R.R. Tolkien", UnitsSold: 30000000, Price: 20},
	}

	mockProvider := &mockBooksProvider{books: expectedBooks}
	repo := NewBooksRepository(mockProvider)

	ctx := context.Background()
	books, err := repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedBooks, books)
}

func TestBooksRepository_GetAll_EmptyResult(t *testing.T) {
	mockProvider := &mockBooksProvider{books: []models.Book{}}
	repo := NewBooksRepository(mockProvider)

	ctx := context.Background()
	books, err := repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.Empty(t, books)
}

func TestBooksRepository_GetAll_ContextPropagation(t *testing.T) {
	mockProvider := &mockBooksProvider{books: []models.Book{}}
	repo := NewBooksRepository(mockProvider)

	ctx := context.WithValue(context.Background(), "test-key", "test-value")
	_, err := repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, mockProvider.ctx)
	assert.Equal(t, ctx, mockProvider.ctx)
}

func TestBooksRepository_GetAll_ProviderReturnsBooks(t *testing.T) {
	mockBooks := []models.Book{
		{ID: 1, Name: "Book 1", Author: "Author 1", UnitsSold: 1000, Price: 10},
		{ID: 2, Name: "Book 2", Author: "Author 2", UnitsSold: 2000, Price: 20},
		{ID: 3, Name: "Book 3", Author: "Author 3", UnitsSold: 3000, Price: 30},
	}

	mockProvider := &mockBooksProvider{books: mockBooks}
	repo := NewBooksRepository(mockProvider)

	ctx := context.Background()
	books, err := repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, books, 3)
	assert.Equal(t, mockBooks, books)
}
