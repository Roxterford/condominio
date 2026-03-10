package api

import "fmt"

type Error struct {
	Message string `json:"message,omitempty"`
} // @name Error

type Response[Data any] struct {
	ApiVersion string `json:"apiVersion"        example:"1.0"`
	Context    string `json:"context,omitempty"`
	Data       Data   `json:"data"`
	Error      Error  `json:"error,omitempty"`
} // @name ApiResponse

func NewResponse[Data any](data Data, err error) Response[Data] {

	var errorMessage string

	if err != nil {
		errorMessage = err.Error()
	}

	res := Response[Data]{
		ApiVersion: "1.0",
		Data:       data,
		Error: Error{
			Message: errorMessage,
		},
	}

	return res
}

func NewErrorResponse(message string, a ...any) Response[any] {
	return NewResponse[any](nil, fmt.Errorf(message, a...))
}
