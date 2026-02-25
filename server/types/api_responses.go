package types

type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse struct {
	Items      any        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

func NewPaginatedResponse(items any, page, perPage, total int) PaginatedResponse {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 1
	}

	totalPages := 1
	if total > 0 {
		totalPages = (total + perPage - 1) / perPage
	}

	return PaginatedResponse{
		Items: items,
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
