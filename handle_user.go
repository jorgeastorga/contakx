package main

import (
	"net/http"
	"log"
	)

/**
*
*
 */
func userHandlerEdit(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	RenderTemplate(w, r, "users/edit", map[string]interface{}{
		"User": user,
	})
}

/**
*
*
*/
func userHandlerUpdate(w http.ResponseWriter, r *http.Request) {
	currentUser := RequestUser(r)
	email := r.FormValue("email")
	currentPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")

	user, err := UpdateUser(currentUser, email, currentPassword, newPassword)
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "users/edit", map[string]interface{}{
				"Error": err.Error(),
				"User":  user,
			})
			return
		}
		panic(err)
	}

	err = globalUserStore.Save(*currentUser)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/account?flash=User+Updated", http.StatusFound)
}

/********************************************************
*  Registration Handler
*  Used to render the registration view/page
*/
func registrationHandlerGET(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "users/new", nil)
}

/********************************************************
*  Registration Handler
*  Used to save the registration information
*/
func registrationHandlerPOST(w http.ResponseWriter, r *http.Request) {
	//create user from form information
	user, err := NewUser(
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password"))

	//error checking for user created
	if err != nil {
		if IsValidationError(err) {
			log.Println(err.Error())
			RenderTemplate(w, r, "/users/new", map[string]interface{}{
				"Error": err.Error(),
				"User":  user,
			})
			panic(err)
			return
		}
	}

	//(Attemp to)Save the user
	err = globalUserStore.Save(user)
	if err != nil {
		panic(err)
		return
	}

	//Create a session for newly created user
	session := NewSession(w)
	session.UserID = user.ID
	err = globalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}

	//Redirect user to account view/page
	http.Redirect(w, r, "/account?flash=User+created", http.StatusFound)
}
