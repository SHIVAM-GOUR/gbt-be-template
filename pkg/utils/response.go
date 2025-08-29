package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// WriteJSONResponse writes a JSON response
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteSuccessResponse writes a successful JSON response
func WriteSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	WriteJSONResponse(w, statusCode, response)
}

// WriteErrorResponse writes an error JSON response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string, err interface{}) {
	response := APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
	WriteJSONResponse(w, statusCode, response)
}

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// WritePaginatedResponse writes a paginated JSON response
func WritePaginatedResponse(w http.ResponseWriter, statusCode int, message string, data interface{}, total int64, page, limit int) {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	
	pagination := PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	response := APIResponse{
		Success: true,
		Message: message,
		Data:    pagination,
	}
	
	WriteJSONResponse(w, statusCode, response)
}
