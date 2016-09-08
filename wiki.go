package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  //"log"
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

func viewHandler(w http.ResponseWriter, r *http.Request){
  title := r.URL.Path[len("/view/"):]
  p,_ := loadPage(title)
  fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func handler(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "Hi there, %s!", r.URL.Path[1:])
}

func main() {
  p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
  p1.save()
  p2,_ := loadPage("TestPage")
  fmt.Println(string(p2.Body))

  port := os.Getenv("PORT")

  if port == "" {
    log.Fatal("$PORT must be set")
    }

    http.HandleFunc("/view/", viewHandler)
    http.ListenAndServe(":" + port, nil)
    //http.ListenAndServe(":8080", nil)
}
