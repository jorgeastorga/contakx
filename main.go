/**
* Context
*
* Main application executable
*/
package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	_ "database/sql"
)

/********************************************************
* Main function to initiate the application
*/
func main() {

	//Route Registration
	unauthenticatedRouter := mux.NewRouter()

	/* Static File Server */
	unauthenticatedRouter.PathPrefix("/assets").Handler(
	http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	notFound := new(NotFound)
	unauthenticatedRouter.NotFoundHandler = notFound
	unauthenticatedRouter.HandleFunc("/", indexHandler)
	unauthenticatedRouter.HandleFunc("/view", viewHandler)
	unauthenticatedRouter.HandleFunc("/edit", editHandler)
	unauthenticatedRouter.HandleFunc("/about", aboutHandler)
	unauthenticatedRouter.HandleFunc("/contact", contactHandler)

	unauthenticatedRouter.HandleFunc("/register", registrationHandlerGET).Methods("GET")
	unauthenticatedRouter.HandleFunc("/register", registrationHandlerPOST).Methods("POST")

	unauthenticatedRouter.HandleFunc("/login", loginSessionHandlerNew).Methods("GET")
	unauthenticatedRouter.HandleFunc("/login", loginSessionHandlerCreate).Methods("POST")


	// Registration of secure routes
	secureRouter := mux.NewRouter()
	secureRouter.HandleFunc("/sign-out", loginSessionHandlerDestroy).Methods("GET")
	secureRouter.HandleFunc("/account", userHandlerEdit).Methods("GET")
	secureRouter.HandleFunc("/account", userHandlerUpdate).Methods("POST")

	//Setup middleware (chain model)
	middleWare := MiddleWare{}
	middleWare.Add(unauthenticatedRouter)
	middleWare.Add(http.HandlerFunc(RequireLogin))
	middleWare.Add(secureRouter)

	//Server startup
	port := os.Getenv("PORT")

	if port == "" {
		//Fix this when deploying a good release of the production application
		//log.Fatal("$PORT must be set")
		log.Println("Port not set, using 8080")
		log.Fatal(http.ListenAndServe(":8080", middleWare))
	} else {
		//Used for Heroku
		log.Fatal(http.ListenAndServe(":"+port, middleWare))
	}
}
