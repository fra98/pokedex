package httperror

// HTTPError represents an HTTP error.
type HTTPError struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"statusCode"`
}

func (e HTTPError) Error() string {
	return "message: " + e.Message
}

// NewHTTPError returns a new HTTP error.
func NewHTTPError(message string, statusCode int) HTTPError {
	return HTTPError{
		Message:    message,
		StatusCode: statusCode,
	}
}
