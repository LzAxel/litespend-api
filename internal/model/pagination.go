package model

type PaginationParams struct {
	Page  int
	Limit int
}

func NewPaginationParams(page, limit int) PaginationParams {
	p := PaginationParams{
		Page:  page,
		Limit: limit,
	}
	p.Validate()
	return p
}

func (p *PaginationParams) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}

func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

func NewPaginationMeta(total int, params PaginationParams) PaginationMeta {
	totalPages := 0
	if params.Limit > 0 {
		totalPages = (total + params.Limit - 1) / params.Limit
	}
	return PaginationMeta{
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}
}

type PaginatedResponse[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

func NewPaginatedResponse[T any](data []T, total int, params PaginationParams) PaginatedResponse[T] {
	return PaginatedResponse[T]{
		Data: data,
		Meta: NewPaginationMeta(total, params),
	}
}
