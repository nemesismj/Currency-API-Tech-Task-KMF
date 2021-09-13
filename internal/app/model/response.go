package model
// ErrorResponse struct
type ErrorResponse struct {
	Message string `json:"error"`
}
// SuccessResponse struct
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    []RCurrency `json:"data,omitempty"`
}
