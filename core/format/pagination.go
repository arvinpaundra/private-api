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

// CalculateOffset calculates the offset for pagination based on page and perPage
func CalculateOffset(page, perPage int) int {
	validPage := getPage(page)
	validPerPage := getLimit(perPage)
	return (validPage - 1) * validPerPage
}

// ValidatePage returns a valid page number (minimum 1)
func ValidatePage(page int) int {
	return getPage(page)
}

// ValidatePerPage returns a valid perPage number (minimum 1, default 10)
func ValidatePerPage(perPage int) int {
	if perPage < 1 {
		return 10
	}
	return perPage
}
