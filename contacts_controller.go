package main

import (
	"log"
	"net/http"
)

func contactHandlerView(w http.ResponseWriter, r *http.Request){
	log.Println("Called the Contact Handler")
	contacts := getAllContacts()
	RenderTemplate(w, r, "contacts/view", map[string]interface{}{
		"Contact1": contacts[0],
		"Contacts": contacts,
	})
}


func getAllContacts() []Contact{
	var contacts []Contact
	if contacts == nil {
		log.Println("contats is nil")
	} else {
		log.Println("contacts is not nil")
	}
	AppDB.Find(&contacts)
	return contacts
}
