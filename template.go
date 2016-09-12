package main

import (
  "fmt"
  "html/template"
  "net/http"
  "log"
  "bytes"
)

var layoutFuncs = template.FuncMap{
  "yield": func() (string, error){
    return "", fmt.Errorf("yield called inappropriately")
  },
}

var layout = template.Must(
  template.New("layout.html").
  Funcs(layoutFuncs).
  ParseFiles("templates/layout.html"),)

var templates = template.Must(
  template.New("t").ParseGlob("templates/**/*.html"))

var errorTemplate = `
<html>
<body>
<h1>Error rendering template %s</h1>
</body>
</html>
`

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string,
  data interface{}){

  funcs := template.FuncMap{
    "yield": func() (template.HTML, error){
      buf := bytes.NewBuffer(nil)
      err := templates.ExecuteTemplate(buf, name, data)
      return template.HTML(buf.String()), err
    },
  }

  layoutClone, _ := layout.Clone()
    layoutClone.Funcs(funcs)
    err := layoutClone.Execute(w, data)


  //err := templates.ExecuteTemplate(w, name, data)
  if err != nil {
    http.Error(
      w,
      fmt.Sprintf(errorTemplate, name, err),
      http.StatusInternalServerError,
    )
  }



  log.Println("Render template called successfully")
}
