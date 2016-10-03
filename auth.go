package main

import(
  "net/http"
)

func AuthenticateRequest(w http.ResponseWriter, r *http.Request){
  //log.Println("AuthenticateRequest is getting called")
  //Redirect the user to login if they are not authenticated
  authenticated := false
  if !authenticated {
    http.Redirect(w, r, "/register", http.StatusFound)
  }
}
