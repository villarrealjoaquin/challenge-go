package services

import (
	"context"
	"slices"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/repositories"
)

type MetricsService interface {
	GetMetrics(ctx context.Context, author string) (*MetricsResponse, error)
}

type metricsService struct {
	booksRepo repositories.BooksRepository
}

type MetricsResponse struct {
	Books                []models.Book `json:"books"`
	MeanUnitsSold        uint          `json:"mean_units_sold"`
	CheapestBook         string        `json:"cheapest_book"`
	BooksWrittenByAuthor uint          `json:"books_written_by_author"`
}

func NewMetricsService(booksRepo repositories.BooksRepository) MetricsService {
	return &metricsService{
		booksRepo: booksRepo,
	}
}

func (s *metricsService) GetMetrics(ctx context.Context, author string) (*MetricsResponse, error) {
	books, err := s.booksRepo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return &MetricsResponse{
			Books:                []models.Book{},
			MeanUnitsSold:        0,
			CheapestBook:         "",
			BooksWrittenByAuthor: 0,
		}, nil
	}

	meanUnitsSold := meanUnitsSold(books)
	cheapestBook := cheapestBook(books)
	booksWrittenByAuthor := booksWrittenByAuthor(books, author)

	return &MetricsResponse{
		Books:                books,
		MeanUnitsSold:        meanUnitsSold,
		CheapestBook:         cheapestBook.Name,
		BooksWrittenByAuthor: booksWrittenByAuthor,
	}, nil
}

func meanUnitsSold(books []models.Book) uint {
	if len(books) == 0 {
		return 0
	}
	var sum uint
	for _, book := range books {
		sum += book.UnitsSold
	}
	return sum / uint(len(books))
}

func cheapestBook(books []models.Book) models.Book {
	return slices.MinFunc(books, func(a, b models.Book) int {
		return int(a.Price - b.Price)
	})
}

func booksWrittenByAuthor(books []models.Book, author string) uint {
	var count uint
	for _, book := range books {
		if book.Author == author {
			count++
		}
	}
	return count
}
