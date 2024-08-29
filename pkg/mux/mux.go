package mux

import "net/http"

func CreateHandler(h http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return h
}
