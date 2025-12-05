package response

type APIError struct {
	StatusCode int
	Message    string
	Data       interface{}
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(status int, msg string, data interface{}) APIError {
	return APIError{
		StatusCode: status,
		Message:    msg,
		Data:       data,
	}
}