package main

import (
  "io/ioutil"
  "net/http"
  "os"
  "log"
  "html/template"
)

type Page struct {
  Title string
  Body []byte
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

func handler(w http.ResponseWriter, r *http.Request){
  //fmt.Fprintf(w, "Hi there, %s!", r.URL.Path[1:])
}

/********************************************************
* Main function to initiate the application
*/
func main() {

  //Route Registration
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
    RenderTemplate(w, r, "index/home", nil)
  })

  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/about", aboutHandler)
  http.HandleFunc("/contact", contactHandler)
  //http.HandleFunc("/save/", saveHandler)


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
