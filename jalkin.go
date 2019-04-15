package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("public_html/copic_markers.html"))

//Page sdas
type Page struct {
	Title string
	Body  []byte
	Data  template.HTML
}

type marker struct {
	Name      string
	ColorCode string
	Acquired  bool
}

func main() {
	http.HandleFunc("/copic_markers/", copicHandler)
	// http.HandleFunc("/edit/", editHandler)
	// http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func copicHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/view/"):]
	p := Page{
		Title: "Copic markers",
		Body:  []byte("test"),
	}
	markers := getAllMarkers()
	p.Data = template.HTML(markers)
	renderTemplate(w, "copic_markers", &p)

}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "public_html/" + title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getAllMarkers() string {
	allMarkers := []marker{
		{"Cobalt Blue", "0047ab", true},
		{"Sap Green", "507d2a", true},
		{"Fairy Red", "ffcccc", true},
	}
	htmlString := ""
	for _, copics := range allMarkers {
		htmlString = fmt.Sprintf("%s\n<div style=\"background-color: #%s ; width: 150px; padding: 10px; border: 1px solid green;\"><li>%s</li></div>", htmlString, copics.ColorCode, copics.Name)

	}

	return htmlString
}
