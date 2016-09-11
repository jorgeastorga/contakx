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


  /*TODO:Stop loading a text file since we're changin this logic
  title := r.URL.Path[len("/view/"):]
  p, err:= loadPage(title)

  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }*/

  renderTemplate(w, "view", nil)
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

/**
* Main function to initiate the application
*/
func main() {

  //Route Registration
  http.HandleFunc("/", viewHandler)
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
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
