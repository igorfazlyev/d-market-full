package handlers

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationQuery represents pagination parameters
type PaginationQuery struct {
	Page    int `form:"page,default=1"`
	PerPage int `form:"per_page,default=20"`
}

// DateRangeQuery represents date range parameters
type DateRangeQuery struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
