/**
* Context
*
* Main application executable
*
*
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
* View Handler
*/
func viewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("viewHandler was called")
}

/********************************************************
* About Handler
*/
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index/about", nil)
}

/********************************************************
* Edit Handler
 */
func editHandler(w http.ResponseWriter, r *http.Request) {
	//Saving for sample code
	/*title := r.URL.Path[len("/edit/"):]
	  p, err := loadPage(title)
	  if err != nil {
	  p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)*/
}

/********************************************************
*  Save Handler
 */
func saveHandler(w http.ResponseWriter, r *http.Request) {

}

/********************************************************
*  Contact Handler
 */
func contactHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index/contact", nil)
}

/********************************************************
*  Index Handler
 */
func indexHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index/home", nil)
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

/********************************************************
*  Gorilla Mux Page Handler
 */
func pageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("pages handler worked")
	vars := mux.Vars(r)
	pageID := vars["id"]

	log.Println(pageID)
}

func testingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("testing handler was called")
	RenderTemplate(w, r, "index/about", nil)
}

/*****
* This NotFound structure is used so that the mux not found handler
* passes along the request to subsequent handlers when a
* URL pattern does not match
 */
type NotFound struct{}

func (n *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

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
