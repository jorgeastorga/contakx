package main

import (
	"net/http"
	"net/url"
	"time"
)

const (
	//Keep users logged in for 3 days
	sessionLength     = 24 * 3 * time.Hour
	sessionCookieName = "ContakxSession"
	sessionIDLength   = 20
)

type Session struct {
	ID     string
	UserID string
	Expiry time.Time
}

func NewSession(w http.ResponseWriter) *Session {
	expiry := time.Now().Add(sessionLength)

	session := &Session{
		ID:     GenerateID("sess", sessionIDLength),
		Expiry: expiry,
	}

	cookie := http.Cookie{
		Name:    sessionCookieName,
		Value:   session.ID,
		Expires: expiry,
	}

	http.SetCookie(w, &cookie)
	return session
}

/**
* Input: HTTP Request
* Output:
*       - Session object when session exists in request (in cookie and not expired)
*       - nil if session is not found in request (cookie) or expired
* Description: the function retrieves the session from a cookie
*              and determines if the session exist in the system.
 */
func RequestSession(r *http.Request) *Session {

	//retrieve the cookie
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil
	}

	//retrieve the session
	session, err := globalSessionStore.Find(cookie.Value)
	if err != nil {
		panic(err)
	}

	if session == nil {
		return nil
	}

	//check if the session has expired
	if session.Expired() {
		globalSessionStore.Delete(session)
		return nil
	}
	return session
}

/**
* Input: N/A (no input - operates on internal data)
* Output:
*        True if session HAS expired
*        False if session is NOT expired
* Description: Function used to check if a session has expired
*
 */
func (session *Session) Expired() bool {
	return session.Expiry.Before(time.Now())
}

/***
* Input: http.Request (object)
* Output:
*        - User (object) when user is found
*        - nil when user not found
* Description: looks up a user assuming the request contains appropriate
*              cookie information
*
 */
func RequestUser(r *http.Request) *User {
	session := RequestSession(r)
	if session == nil || session.UserID == "" {
		return nil
	}

	user, err := globalUserStore.Find(session.UserID)
	if err != nil {
		panic(err)
	}
	return user
}

/***
* Input: http Request, http ResponseWriter
* Output:
* Description: Handler used to
 */
func RequireLogin(w http.ResponseWriter, r *http.Request) {
	//Let the request pass if we've got a user
	if RequestUser(r) != nil {
		return // user has been found
	}

	//If the user is not found, require the user to log in
	query := url.Values{}
	//Retain the desired url to redirect after login
	query.Add("next", url.QueryEscape(r.URL.String()))
	http.Redirect(w, r, "/login?"+query.Encode(), http.StatusFound)
}

/**
*
*
 */
func FindOrCreateSession(w http.ResponseWriter, r *http.Request) *Session {
	session := RequestSession(r)
	if session == nil {
		session = NewSession(w)
	}

	return session
}
