package main

import "net/http"

func loginSessionHandlerDestroy(w http.ResponseWriter, r *http.Request) {
	session := RequestSession(r)
	if session != nil {
		err := globalSessionStore.Delete(session)
		if err != nil {
			panic(err)
		}
	}
	RenderTemplate(w, r, "sessions/destroy", nil)
}
