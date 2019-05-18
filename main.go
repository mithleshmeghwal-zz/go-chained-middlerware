package main

import (
	"net/http"
	"time"
	"os"
	"log"
	"fmt"
)

type Logger struct {
	handler http.Handler
}

func (l *Logger)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	end := time.Since(start)
	log.Printf("%s %s %v", r.Method, r.URL.Path, end)
}

func NewLogger(h http.Handler) *Logger {
	return &Logger{h}
}

type ResponseHeader struct {
	handler http.Handler
	headerName string
	headerValue string
}

func NewResponseHeader(h http.Handler, headerName, headerValue string) *ResponseHeader {
	return &ResponseHeader{ h, headerName, headerValue }
}

func (rh *ResponseHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(rh.headerName, rh.headerValue)
	rh.handler.ServeHTTP(w, r)
}


func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, world"))
}

func CurrentTimeHandler(w http.ResponseWriter, r *http.Request) {
	curTime := time.Now().Format(time.Kitchen)
	w.Write([]byte(fmt.Sprintf("the current time is %v", curTime)))
}

func main() {
	addr := os.Getenv("ADDR")
	
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/hello", HelloHandler)
	mux.HandleFunc("/v1/time", CurrentTimeHandler)

	wrappedMux := NewLogger(NewResponseHeader(mux, "X-My-Header", "my header value"))

	log.Printf("server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(":"+addr, wrappedMux))
}