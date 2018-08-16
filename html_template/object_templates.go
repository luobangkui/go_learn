package main

import (
	"html/template"
	"bytes"
	"net/http"
)
var t *template.Template
var qc template.HTML
func init() {
	t = template.Must(template.ParseFiles("./templates/index.html", "./templates/quote.html"))
}
type Page struct {
	Title string
	Content template.HTML
}
type Quote struct {
	Quote, Person string
}

func main() {
	q := &Quote{
		Quote: `You keep using that word. I do not think
 it means what you think it means.`,
		Person: "Inigo Montoya",
	}
	var b bytes.Buffer
	t.ExecuteTemplate(&b, "quote.html", q)
	qc = template.HTML(b.String())
	http.HandleFunc("/", diaplayPage)
	http.ListenAndServe(":8081", nil)
}

func diaplayPage(w http.ResponseWriter, r *http.Request) {
	p := &Page{
		Title: "A User",
		Content: qc,
	}
	t.ExecuteTemplate(w, "index.html", p)
}