package types

import "fmt"

// Pagination represents pagination parameters as a value object
type Pagination struct {
	limit  int
	offset int
}

// NewPagination creates a new Pagination value object
func NewPagination(limit, offset int) (Pagination, error) {
	if limit < 1 {
		return Pagination{}, fmt.Errorf("limit must be at least 1")
	}

	if limit > 1000 {
		return Pagination{}, fmt.Errorf("limit cannot exceed 1000")
	}

	if offset < 0 {
		return Pagination{}, fmt.Errorf("offset cannot be negative")
	}

	return Pagination{
		limit:  limit,
		offset: offset,
	}, nil
}

// DefaultPagination creates pagination with default values
func DefaultPagination() Pagination {
	return Pagination{
		limit:  50,
		offset: 0,
	}
}

// Limit returns the limit value
func (p Pagination) Limit() int {
	return p.limit
}

// Offset returns the offset value
func (p Pagination) Offset() int {
	return p.offset
}

// Page returns the current page number (1-based)
func (p Pagination) Page() int {
	return (p.offset / p.limit) + 1
}

// HasNext checks if there are more pages given a total count
func (p Pagination) HasNext(totalCount int) bool {
	return p.offset+p.limit < totalCount
}

// HasPrevious checks if there are previous pages
func (p Pagination) HasPrevious() bool {
	return p.offset > 0
}

// Next returns pagination for the next page
func (p Pagination) Next() Pagination {
	return Pagination{
		limit:  p.limit,
		offset: p.offset + p.limit,
	}
}

// Previous returns pagination for the previous page
func (p Pagination) Previous() Pagination {
	newOffset := p.offset - p.limit
	if newOffset < 0 {
		newOffset = 0
	}
	return Pagination{
		limit:  p.limit,
		offset: newOffset,
	}
}

// WithLimit returns pagination with a different limit
func (p Pagination) WithLimit(limit int) (Pagination, error) {
	return NewPagination(limit, p.offset)
}

// WithOffset returns pagination with a different offset
func (p Pagination) WithOffset(offset int) (Pagination, error) {
	return NewPagination(p.limit, offset)
}
