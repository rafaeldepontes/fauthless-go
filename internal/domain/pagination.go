package domain

type OffSetPagination[T any] struct {
	Data         []T  `json:"data"`
	CurrentPage  uint `json:"currentPage"`
	TotalPages   uint `json:"totalPages"`
	TotalRecords uint `json:"totalRecords"`
	Size         uint `json:"size"`
}

type CursorPagination[T any] struct {
	Data       []T   `json:"data"`
	Size       int   `json:"size"`
	NextCursor int64 `json:"next_cursor"`
}

type CursorResquest struct {
	Cursor string `json:"cursor"`
}

type CursorHashedPagination[T any] struct {
	NextCursor string `json:"next_cursor"`
}
