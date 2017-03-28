package main

import (
	"log"
	"net/http"
)

func contactHandlerView(w http.ResponseWriter, r *http.Request){
	log.Println("Called the Contact Handler")
	RenderTemplate(w, r, "contacts/view", nil)
}
