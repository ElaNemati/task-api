package models

type TaskQueryParams struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Sort   string `form:"sort"`
	Status Status `form:"status"`
}

type PaginatedResponse struct {
	Tasks      []Task `json:"tasks"`
	TotalCount int64  `json:"total_count"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}