package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var templates = template.Must(template.ParseFiles("public_html/copic_markers.html"))

//Page dsdasd
type Page struct {
	Title           string
	Body            []byte
	Data            template.HTML
	AcquiredMarkers template.HTML
	WantedMarkers   template.HTML
}

type marker struct {
	Name      string
	ColorCode string
	Acquired  bool
}

type htmlLoc struct {
	htmlString string
	htmlLoc    string
}

func main() {
	http.HandleFunc("/copic_markers/", copicHandler)
	http.HandleFunc("/copic_markers/*", singleMarkerHandler)
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
	markerName := r.URL.Query()["name"]
	if markerName != nil {
		p.Data = template.HTML(getOneMarker(markerName[0]).Name)
		renderTemplate(w, "marker_page", &p)
	}
	markers := getAllMarkers()
	p.AcquiredMarkers = template.HTML(markers[0].htmlString)
	p.WantedMarkers = template.HTML(markers[1].htmlString)
	renderTemplate(w, "copic_markers", &p)

}

func singleMarkerHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{
		Title: "Copic markers",
		Body:  []byte("test"),
	}
	renderTemplate(w, "marker_page", &p)
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

func getOneMarker(markerName string) marker {
	allMarkers := []marker{
		{"Cobalt Blue", "0047ab", true},
		{"Sap Green", "507d2a", true},
		{"Fairy Red", "ffcccc", true},
		{"Berry Cool", "ff00ff", false},
	}
	return allMarkers[0]
}

func getAllMarkers() []htmlLoc {
	allMarkers := []marker{
		{"Cobalt Blue", "0047ab", true},
		{"Sap Green", "507d2a", true},
		{"Fairy Red", "ffcccc", true},
		{"Berry Cool", "ff00ff", false},
	}
	var result []htmlLoc
	var htmlAcq string
	var htmlWanted string
	for _, copics := range allMarkers {
		if copics.Acquired {
			urlName := strings.Replace(copics.Name, " ", "_", -1)
			htmlAcq = fmt.Sprintf("%s<a href= \"%s\"><div style=\"background-color: #%s ; width: 150px; padding: 10px; border: 1px solid green;\">%s</div></a>", htmlAcq, urlName, copics.ColorCode, copics.Name)
		} else {
			urlName := strings.Replace(copics.Name, " ", "_", -1)
			htmlWanted = fmt.Sprintf("%s<a href= \"%s\"><div style=\"background-color: #%s ; width: 150px; padding: 10px; border: 1px solid green;\">%s</div></a>", htmlWanted, urlName, copics.ColorCode, copics.Name)
		}
	}
	aquiredMarker := htmlLoc{
		htmlString: htmlAcq,
		htmlLoc:    "AquiredMarkers",
	}
	wantedMarker := htmlLoc{
		htmlString: htmlWanted,
		htmlLoc:    "WantedMarkers",
	}
	result = append(result, aquiredMarker)
	result = append(result, wantedMarker)

	return result
}
