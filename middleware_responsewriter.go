package main

import (
  "net/http"
)

/***********************************
* MiddleWareResponseWriter Structure
*
*/
type MiddleWareResponseWriter struct {
  http.ResponseWriter
  written bool
}

/***********************************
* ----
*/
func NewMiddlewareResponseWriter(w http.ResponseWriter) *MiddleWareResponseWriter {
  return &MiddleWareResponseWriter{
    ResponseWriter: w,
  }
}

/***********************************
* ----
*/
func (w *MiddleWareResponseWriter) Write(bytes []byte) (int, error){
  w.written = true
  return w.ResponseWriter.Write(bytes)
}

/***********************************
* ----
*/
func (w MiddleWareResponseWriter) WriteHeader(code int){
  w.written = true
  w.ResponseWriter.WriteHeader(code)
}
