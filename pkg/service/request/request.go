package request

import (
	"errors"
	"io"
	"net/http"
)

var (
	ErrInvalidRequest = errors.New("[scenicTicket] create request failed")
)

// Get -
func Get(URL string, body io.Reader) *http.Request {
	req, err := http.NewRequest("GET", URL, nil)

	if err != nil {
		return nil
	}

	return req
}

// Post
func Post(URL string, body io.Reader) *http.Request {
	req, err := http.NewRequest("POST", URL, nil)

	if err != nil {
		return nil
	}

	return req
}
