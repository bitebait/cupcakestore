package models

const (
	DefaultPage  = 1
	DefaultLimit = 8
	MaxLimit     = 100
)

type Pagination struct {
	Page  int
	Limit int
	Total int64
}

func NewPagination(page, limit int) *Pagination {
	if page <= 0 {
		page = DefaultPage
	}
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	// Bound offsets so malicious query parameters cannot overflow an int.
	maxPage := int(^uint(0)>>1) / limit
	if page > maxPage {
		page = maxPage
	}
	return &Pagination{
		Page:  page,
		Limit: limit,
	}
}

func (p *Pagination) TotalPages() int {
	if p.Limit <= 0 || p.Total <= 0 {
		return 0
	}
	totalPages := int(p.Total) / p.Limit
	if int(p.Total)%p.Limit != 0 {
		totalPages++
	}

	return totalPages
}
