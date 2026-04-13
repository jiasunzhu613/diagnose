package api

type StatusError struct {
	error
	StatusCode int
	Status string
	ErrorMessage string
}

type AuthorizationError struct {
	error
	StatusCode int
	Status string
	SigninURL string
}
