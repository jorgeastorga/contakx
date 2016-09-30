package main

import (
  "net/http"
  "os"
  "log"
  _"database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/gorilla/mux"
  _"github.com/jinzhu/gorm"
)

/********************************************************
* View Handler
*/
func viewHandler(w http.ResponseWriter, r *http.Request){
  log.Println("viewHandler was called")
}

/********************************************************
* View Handler
*/
func aboutHandler(w http.ResponseWriter, r *http.Request){
  RenderTemplate(w, r, "index/about", nil)
}

/********************************************************
* Edit Handler
*/
func editHandler(w http.ResponseWriter, r *http.Request){
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
func saveHandler(w http.ResponseWriter, r *http.Request){

}

/********************************************************
*  Contact Handler
*/
func contactHandler(w http.ResponseWriter, r *http.Request){
  RenderTemplate(w, r, "index/contact", nil)
}

/********************************************************
*  Index Handler
*/
func indexHandler(w http.ResponseWriter, r *http.Request){
  RenderTemplate(w, r, "index/home", nil)
}

/********************************************************
*  Registration Handler
*/
func registrationHandler(w http.ResponseWriter, r *http.Request){
  RenderTemplate(w, r, "users/new", nil)
}

/********************************************************
*  Gorilla Mux Page Handler
*/
func pageHandler(w http.ResponseWriter, r *http.Request){
  log.Println("pages handler worked")
  vars := mux.Vars(r)
  pageID := vars["id"]

  log.Println(pageID)
}

func testingHandler(w http.ResponseWriter, r *http.Request){
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
  //rtr := mux.NewRouter()

  unauthenticatedRouter := mux.NewRouter()
  unauthenticatedRouter.HandleFunc("/", indexHandler)

  notFound := new(NotFound)
  unauthenticatedRouter.NotFoundHandler = notFound
  unauthenticatedRouter.HandleFunc("/view", viewHandler)
  unauthenticatedRouter.HandleFunc("/edit", editHandler)
  unauthenticatedRouter.HandleFunc("/about", aboutHandler)
  unauthenticatedRouter.HandleFunc("/contact", contactHandler)
  unauthenticatedRouter.HandleFunc("/register", registrationHandler)

  /* Static File Server */
  unauthenticatedRouter.PathPrefix("/assets").Handler(
    http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))


  /*authenticatedRouter := mux.NewRouter()
  authenticatedRouter.HandleFunc("/testing/new", testingHandler)*/

  middleWare := MiddleWare{}
  middleWare.Add(unauthenticatedRouter)
  //middleWare.Add(http.HandlerFunc(AuthenticateRequest))
  //middleWare.Add(authenticatedRouter)

  //Server startup
  port := os.Getenv("PORT")


  if port == "" {
    //Fix this when deploying a good release of the production application
    //log.Fatal("$PORT must be set")
    log.Println("Port not set, using 8080")
    //http.ListenAndServe(":8080", rtr)
    log.Fatal(http.ListenAndServe(":8080", middleWare))
    } else {
      //Used for Heroku
      //http.ListenAndServe(":" + port, rtr)
      log.Fatal(http.ListenAndServe(":" + port, middleWare))
    }
  }
