package main


import (
  "net/http"
  "log"
)

/********************************************************
* View Handler
*/
func viewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("viewHandler was called")
}

/********************************************************
* About Handler
*/
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index/about", nil)
}

/********************************************************
* Edit Handler
 */
func editHandler(w http.ResponseWriter, r *http.Request) {
	//Saving for sample code
	/*title := r.URL.Path[len("/edit/"):]
	  p, err := loadPage(title)
	  if err != nil {
	  p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)*/
}

/********************************************************
*  Contact Handler
 */
func contactHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index/contact", nil)
}

/********************************************************
*  Index Handler
 */
func indexHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index/home", nil)
}

/*****
* This NotFound structure is used so that the mux not found handler
* passes along the request to subsequent handlers when a
* URL pattern does not match
 */
type NotFound struct{}

func (n *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
