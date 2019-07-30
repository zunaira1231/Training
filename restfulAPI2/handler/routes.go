package handler

import (
	"net/http"
	"restfulAPI2/db"
)
func SetUpRouting() *http.ServeMux {
	todoHandler := &todoHandler{
		samples:  &db.Sample{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/samples", todoHandler.GetSamples)

	return mux
}