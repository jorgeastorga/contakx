package main

import (
  "fmt"
  "net/http"
  "os"
  "log"
)

func handler(w http.ResponseWriter, r *http.Request){
  fmt.Fprint(w, "Hi there, %s!", r.URL.Path[1:])
}

func main(){
  http.HandleFunc("/", handler)
  port := os.Getenv("PORT")

  if port == "" {
    log.Fatal("$PORT must be set")
    }

  //http.ListenAndServe(":8080", nil)
  http.ListenAndServe(":" + port, nil)
}
