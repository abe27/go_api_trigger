package models

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code,omitempty"`
	Success bool        `json:"success,omitempty"`
	Error   string      `json:"error,omitempty"`
}
