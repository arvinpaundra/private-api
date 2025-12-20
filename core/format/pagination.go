package format

import "math"

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func getLimit(limit int) int {
	if limit < 1 {
		return 1
	}

	return limit
}

func getPage(page int) int {
	if page < 1 {
		return 1
	}

	return page
}

func getTotalPages(total, limit int) int {
	totalPages := math.Ceil(float64(total) / float64(limit))

	if totalPages < 1 {
		return 1
	}

	return int(totalPages)
}

// NewPagination common pagination technique using limit and offset
func NewPagination(page, perPage, total int) Pagination {
	return Pagination{
		Page:       getPage(page),
		PerPage:    getLimit(perPage),
		TotalPages: getTotalPages(total, perPage),
		Total:      total,
	}
}
