package repositories

import (
	"context"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/providers"
)

type BooksRepository interface {
	GetAll(ctx context.Context) ([]models.Book, error)
}

type booksRepository struct {
	provider providers.BooksProvider
}

func NewBooksRepository(provider providers.BooksProvider) BooksRepository {
	return &booksRepository{
		provider: provider,
	}
}

func (r *booksRepository) GetAll(ctx context.Context) ([]models.Book, error) {
	books := r.provider.GetBooks(ctx)
	return books, nil
}
