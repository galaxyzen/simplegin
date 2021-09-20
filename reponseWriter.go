package simplegin

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
	written bool
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}


func (w *responseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}
