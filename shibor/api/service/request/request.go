package service

import (
	"io"
	"net/http"
)

func Get(URL string, body io.Reader) *http.Request {
	req, err := http.NewRequest("GET", URL, body)
	if err != nil {
		return nil
	}

	return req
}
