package main

import (
	"io"
	"net/http"
)

type ShiborAverage struct {
	TermCode   string  `json:"termCode"`
	ActionCode string  `json:"actionCode"`
	List       []Value `json:"list"`
}

type Value struct {
}

func Get(URL string, body io.Reader) *http.Request {
	req, err := http.NewRequest("GET", URL, body)
	if err != nil {
		return nil
	}

	return req
}
