package model

type Response[T any] struct {
	Status  string `json:"status" extensions:"x-order=0"`
	Message string `json:"message" extensions:"x-order=1"`
	Data    T      `json:"data" extensions:"x-order=2"`
}

func NewResponse[T any](status string, message string, data T) *Response[T] {
	return &Response[T]{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
