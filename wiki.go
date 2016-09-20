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
)

type Page struct {
  Title string
  Body []byte
  Content string
  Date string
}

const (
  DBHost  = "127.0.0.1"
  DBPort  = ":3306"
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
* Main function to initiate the application
*/
func main() {

  //Database connection
  dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBDbase)
  db, err := sql.Open("mysql", dbConn)
  if err != nil {
    log.Println("Couldn't connect!")
    log.Println(err.Error)
  } else {
    log.Println("Connected to the database succesfully")
  }
  database = db


  //Route Registration
  http.HandleFunc("/", indexHandler )
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/about", aboutHandler)
  http.HandleFunc("/contact", contactHandler)
  http.HandleFunc("/register", registrationHandler)


  //Setup the file server to serve assets
  http.Handle("/assets/",
    http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

  //Server startup
  port := os.Getenv("PORT")


  if port == "" {
    //Fix this when deploying a good release of the production application
    //log.Fatal("$PORT must be set")
    log.Println("Port not set, using 8080")
    http.ListenAndServe(":8080", nil)
    } else {
      //Used for Heroku
      http.ListenAndServe(":" + port, nil)
    }
  }
