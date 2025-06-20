package domain

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Object  interface{} `json:"object"`
	Errors  []string    `json:"errors"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Object     interface{} `json:"object"`
	PageNumber int         `json:"page_number"`
	PageSize   int         `json:"page_size"`
	TotalSize  int         `json:"total_size"`
	Errors     []string    `json:"errors"`
}
