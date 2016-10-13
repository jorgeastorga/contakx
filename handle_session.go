package main

import (
	"log"
	"net/http"
)

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

/**
* Handles the user sign in process
*
 */
func loginSessionHandlerNew(w http.ResponseWriter, r *http.Request) {
	next := r.URL.Query().Get("next")
	RenderTemplate(w, r, "sessions/new", map[string]interface{}{
		"Next": next,
	})
}

func loginSessionHandlerCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("the login create function is being called")
	//Extract values from the form
	username := r.FormValue("username")
	password := r.FormValue("password")
	next := r.FormValue("next")

	//Attempt to find the user for the supplied username + password
	user, err := FindUser(username, password)
	if err != nil {
		//If there is a problem with fetching the user, display an error in login page
		if IsValidationError(err) {
			RenderTemplate(w, r, "/sessions/new", map[string]interface{}{
				"Error": err,
				"User":  user,
				"Next":  next,
			})
			return
		}
		panic(err)
	}

	//If user found, we find an existing session or create a new one
	session := FindOrCreateSession(w, r)
	session.UserID = user.ID
	err = globalSessionStore.Save(session)

	if err != nil {
		panic(err)
	}

	if next == "" {
		next = "/"
	}

	//Redirect to the page the user wanted to go
	http.Redirect(w, r, next+"?flash=Signed+in", http.StatusFound)

}
