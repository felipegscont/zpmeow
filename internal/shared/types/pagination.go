package types

import "fmt"

type Pagination struct {
	limit  int
	offset int
}

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

func DefaultPagination() Pagination {
	return Pagination{
		limit:  50,
		offset: 0,
	}
}

func (p Pagination) Limit() int {
	return p.limit
}

func (p Pagination) Offset() int {
	return p.offset
}

func (p Pagination) Page() int {
	return (p.offset / p.limit) + 1
}

func (p Pagination) HasNext(totalCount int) bool {
	return p.offset+p.limit < totalCount
}

func (p Pagination) HasPrevious() bool {
	return p.offset > 0
}

func (p Pagination) Next() Pagination {
	return Pagination{
		limit:  p.limit,
		offset: p.offset + p.limit,
	}
}

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

func (p Pagination) WithLimit(limit int) (Pagination, error) {
	return NewPagination(limit, p.offset)
}

func (p Pagination) WithOffset(offset int) (Pagination, error) {
	return NewPagination(p.limit, offset)
}
