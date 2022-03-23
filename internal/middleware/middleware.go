package middleware

import (
	"log"
	"net/http"
	"time"
)

//Logger is a middleware handler for logging incoming request
type Logger struct {
	handler http.Handler
}

//ResponseHeader is a middleware handler that adds a header to the response
type ResponseHeader struct {
	handler     http.Handler
	headerName  string
	headerValue string
}

//ServeHTTP handles  and logging the request details
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s from %v", r.Method, r.URL.Path, r.RemoteAddr)

	start := time.Now()

	l.handler.ServeHTTP(w, r)

	log.Printf("%s %s duration: %v", r.Method, r.URL.Path, time.Since(start))
}

//NewLogger constructs a new Logger middleware handler
func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}

//NewResponseHeader constructs a new ResponseHeader middleware handler
func NewResponseHeader(handlerToWrap http.Handler, headerName string, headerValue string) *ResponseHeader {
	return &ResponseHeader{handlerToWrap, headerName, headerValue}
}

//ServeHTTP handles the request by adding the response header
func (rh *ResponseHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(rh.headerName, rh.headerValue)

	rh.handler.ServeHTTP(w, r)
}
