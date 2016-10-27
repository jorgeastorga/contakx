package main

import (
	"net/http"
)

/***********************************
* MiddleWareResponseWriter Structure
*
 */
type MiddleWare []http.Handler

/*******************************
* Add
* Description: adds a handler to the Middleware
 */
func (m *MiddleWare) Add(handler http.Handler) {
	*m = append(*m, handler)
}

/*******************************
* ServeHTTP method
 */
func (m MiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Process the middleware
	//Wrap the supplies ResponseWriter
	mw := NewMiddlewareResponseWriter(w)

	//Loop through all the registered handlers
	for _, handler := range m {
		//Call the handler with our MiddleWareResponseWriter
		handler.ServeHTTP(mw, r)

		//If there was a write, stop processing
		if mw.written {
			return
		}
	}
	//If no handlers wrote to the response, it's a 404
	http.NotFound(w, r)
}
