package pagination

import (
	"github.com/rafaeldepontes/auth-go/internal/domain"
)

// NewOffSetPagination accepts a generic T type, a slice of any data, a current number of pages,
// a total number of records and a size and it will return a OffSetPagination.
func NewOffSetPagination[T any](data []T, currentPage, totalRecords, size uint) *domain.OffSetPagination[T] {
	totalPages := totalRecords / size
	if totalPages <= 0 {
		totalPages = 1
	}

	return &domain.OffSetPagination[T]{
		Data:         data,
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		Size:         size,
	}
}

// NewCursorPagination accepts a generic T type, a slice of any data, a size of records per page and
// the next page being a pointer to the next id in the database and it will return a CursorPagination.
func NewCursorPagination[T any](data []T, size int, nextCursor int64) *domain.CursorPagination[T] {
	return &domain.CursorPagination[T]{
		Data:       data,
		Size:       size,
		NextCursor: nextCursor,
	}
}
