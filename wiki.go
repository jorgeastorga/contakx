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

/*
 * View Handler
 */
func viewHandler(w http.ResponseWriter, r *http.Request){
  title := r.URL.Path[len("/view/"):]
  p, err:= loadPage(title)

  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }

  renderTemplate(w, "view", p)
}

/*
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

/*
 *  Save Handler
 */
func saveHandler(w http.ResponseWriter, r *http.Request){

}

func handler(w http.ResponseWriter, r *http.Request){
  //fmt.Fprintf(w, "Hi there, %s!", r.URL.Path[1:])
}

func main() {

  //Route Registration
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  //http.HandleFunc("/save/", saveHandler)

  //Server startup
  port := os.Getenv("PORT")

  if port == "" {
    log.Fatal("$PORT must be set")
  }
  http.ListenAndServe(":" + port, nil)

  //http.ListenAndServe(":8080", nil)
}
