package model

type ErrorResponse struct {
	Title string `json:"title"`
	Error string `json:"error"`
}

type SuccessResponse struct {
	Title string `json:"title"`
}
