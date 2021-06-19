package entity

type ResponseError struct {
	Message    string
	StatusCode int
	Error      error
}
