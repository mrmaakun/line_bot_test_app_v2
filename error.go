package main

import ()

type APIError struct {
	Code     int
	Response string
}

func (e *APIError) Error() string {
	return e.Response
}
