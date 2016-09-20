package main

import (
  "io/ioutil"
  "net/http"
  "os"
  "log"
  "html/template"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "fmt"
  "github.com/gorilla/mux"
)

type Page struct {
  Title string
  Body []byte
  Content string
  Date string
}

const (
  DBHost  = "127.0.0.1"
  DBPort  = "3306"
  DBUser  = "root"
  DBPass  = "testing123"
  DBDbase = "contakx"
)

var database *sql.DB

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
  thisPage := Page{}
  fmt.Println(pageID)
  err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE id=?", pageID).Scan(&thisPage.Title, &thisPage.Content, &thisPage.Date)
  if err != nil {

    log.Println("Couldn't get page: +pageID")
    log.Println(err.Error)
    log.Fatal(err)
  }
  html := `<html><head><title>` + thisPage.Title + `</title></head><body><h1>` + thisPage.Title + `</h1><div>` + thisPage.Content + `</div></body></html>`
  fmt.Fprintln(w, html)
}

/********************************************************
* Main function to initiate the application
*/
func main() {

  //Database connection
  dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUser, DBPass, DBHost, DBPort, DBDbase)
  //dbConn := "root:testing123@tcp(localhost:3306)/contakx"
  db, err := sql.Open("mysql", dbConn)
  if err != nil {
    log.Println("Couldn't connect!")
    log.Println(err.Error)
  } else {
    log.Println("Connected to the database succesfully")
  }
  database = db


  //Route Registration
  rtr := mux.NewRouter()

  /*rtr.HandleFunc("/pages/{id:[0-9]+}",
    pageHandler)*/
    rtr.HandleFunc("/pages/{id:[0-9]+}",
      servePage)
  rtr.HandleFunc("/", indexHandler)
  rtr.HandleFunc("/view", viewHandler)
  rtr.HandleFunc("/edit", editHandler)
  rtr.HandleFunc("/about", aboutHandler)
  rtr.HandleFunc("/contact", contactHandler)
  rtr.HandleFunc("/register", registrationHandler)

  http.Handle("/", rtr)

  rtr.PathPrefix("/assets").Handler(
    http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

  //Server startup
  port := os.Getenv("PORT")


  if port == "" {
    //Fix this when deploying a good release of the production application
    //log.Fatal("$PORT must be set")
    log.Println("Port not set, using 8080")
    http.ListenAndServe(":8080", rtr)
    } else {
      //Used for Heroku
      http.ListenAndServe(":" + port, rtr)
    }
  }
