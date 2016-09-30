package main

import (
  "io/ioutil"
  "net/http"
  "os"
  "log"
  "html/template"
  _"database/sql"
  _ "github.com/go-sql-driver/mysql"
  "fmt"
  "github.com/gorilla/mux"
  _"github.com/jinzhu/gorm"
)

type Page struct {
  Title string
  Body []byte
  Content string
  Date string
}



func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body,err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }

  return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
  t,_ := template.ParseFiles("layout.html", tmpl + ".html")
  t.ExecuteTemplate(w, "layout", p)
}

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
  title := r.URL.Path[len("/edit/"):]
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
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


func servePage(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  pageID := vars["id"]
  //thisPage := Page{}
  fmt.Println(pageID)
  /*err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE id=?", pageID).Scan(&thisPage.Title, &thisPage.Content, &thisPage.Date)
  if err != nil {

    log.Println("Couldn't get page: +pageID")
    log.Println(err.Error)
    log.Fatal(err)
  }
  html := `<html><head><title>` + thisPage.Title + `</title></head><body><h1>` + thisPage.Title + `</h1><div>` + thisPage.Content + `</div></body></html>`*/
  html := `<html><head><title>Jorge</title></div></body></html>`
  fmt.Fprintln(w, html)
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

  unauthenticatedRouter.HandleFunc("/pages/{id:[0-9]+}",
      servePage)
  unauthenticatedRouter.HandleFunc("/view", viewHandler)
  unauthenticatedRouter.HandleFunc("/edit", editHandler)
  unauthenticatedRouter.HandleFunc("/about", aboutHandler)
  unauthenticatedRouter.HandleFunc("/contact", contactHandler)
  unauthenticatedRouter.HandleFunc("/register", registrationHandler)


  unauthenticatedRouter.PathPrefix("/assets").Handler(
    http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

  authenticatedRouter := mux.NewRouter()
  authenticatedRouter.HandleFunc("/testing/new", testingHandler)

  middleWare := MiddleWare{}
  middleWare.Add(unauthenticatedRouter)
  middleWare.Add(http.HandlerFunc(AuthenticateRequest))
  middleWare.Add(authenticatedRouter)

  //Server startup
  port := os.Getenv("PORT")


  if port == "" {
    //Fix this when deploying a good release of the production application
    //log.Fatal("$PORT must be set")
    log.Println("Port not set, using 8080")
    //http.ListenAndServe(":8080", rtr)
    http.ListenAndServe(":8080", middleWare)
    } else {
      //Used for Heroku
      //http.ListenAndServe(":" + port, rtr)
      http.ListenAndServe(":" + port, middleWare)
    }
  }
