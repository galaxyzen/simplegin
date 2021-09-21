package simplegin

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
	written bool
}

func (w *responseWriter) WriteHeader(statusCode int) {
	if !w.written {
		w.StatusCode = statusCode
		w.written = true
	}
	w.ResponseWriter.WriteHeader(statusCode)
}


func (w *responseWriter) Write(data []byte) (int, error) {
	if !w.written {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}
