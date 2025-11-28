package pagination

import "github.com/rafaeldepontes/auth-go/internal/domain"

// NewPagination accepts a generic T type, a slice of any data, a current number of pages,
// a total number of records and a size and it will return a Pagination.
func NewPagination[T any](data []T, currentPage, totalRecords, size uint) *domain.Pagination[T] {
	totalPages := totalRecords / size
	if totalPages <= 0 {
		totalPages = 1
	}

	return &domain.Pagination[T]{
		Data:         data,
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		Size:         size,
	}
}
