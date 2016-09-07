package main

import (
  "fmt"
  "net/http"
  "os"
  "log"
)

func handler(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "Hi there, %s!", r.URL.Path[1:])
}

func main(){

  //Routing: mapping the handler function to a request path
  //This is the base funcitonality for web frameworks. Map a request to
  //a function.
  http.HandleFunc("/", handler)


  port := os.Getenv("PORT")

  if port == "" {
    log.Fatal("$PORT must be set")
    }

  //http.ListenAndServe(":8080", nil)
  http.ListenAndServe(":" + port, nil)
}
