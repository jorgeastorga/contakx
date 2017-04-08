/**
* Context
*
* Main application executable
 */
package main

import (
	_ "database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"os"
	"fmt"
)


/***********************
* DB Connection Details
 */

const (
	DBHost = "127.0.0.1"
	DBPort = "5432"
	DBUser = "root"
	DBPass = "testing123"
	DBDatabase = "contakx"
)

//Database identified for gorm.DB
var AppDB *gorm.DB


/********************************************************
* Main function to initiate the application
 */
func main() {

	var err error

	//Database Connection Setup: PostgreSQL
	dbConnection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		DBHost,
		DBPort,
		DBUser,
		DBDatabase,
		DBPass)

	AppDB, err = gorm.Open("postgres", dbConnection)

	if err != nil{
		log.Println("Failed to connect to the database", err.Error())
	} else {
		log.Println("DB Connection: connected to the database successfully")
	}


	//Setup Database Tables
	AppDB.AutoMigrate(&Contact{})


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

	secureRouter.HandleFunc("/contacts", contactHandlerView).Methods("GET")

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

//TODO: Delete this function and put into some kind of automated tests
/* Test function to create a contact */
func createContact(){
	newContact := &Contact{
		FirstName: "Esteban",
		LastName: "Renteria",
		Address1: "3230 Stuart Street",
		Address2: "Street",
		City: "Oakland",
		State: "CA",
		ZipCode: 94602,
		Company: "Contaks Inc.",
		Email: "jeastorga@gmail.com",
		Phone: "(123)123-1235", }
	newContact.create()
}