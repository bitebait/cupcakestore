package models

const (
	DefaultPage  = 1
	DefaultLimit = 8
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
	return &Pagination{
		Page:  page,
		Limit: limit,
	}
}

func (p *Pagination) TotalPages() int {
	totalPages := int(p.Total) / p.Limit
	if int(p.Total)%p.Limit != 0 {
		totalPages++
	}

	return totalPages
}
