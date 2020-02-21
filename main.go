package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Page struct {
	Title string
	Body  []byte
}

// save To save a page.
// usage: page1.save()
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// loadPage To load a page with given title
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	fmt.Println("saveHandler(): " + string(p.Body))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/create/"):]
	fmt.Fprintf(w, "<h1>Creating: %s</h1>"+
        "<form action=\"/save/%s\" method=\"POST\">"+
        "<textarea name=\"body\"></textarea><br>"+
        "<input type=\"submit\" value=\"Save\">"+
        "</form>",
        title, title)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/read/"):]
	page, err := loadPage(title)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", page.Title, page.Body)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/update/"):]
	p, err := loadPage(title)
	if err != nil {
		fmt.Fprintf(w, "<h1>Unable to open %s</h1>", title)
		return
	}
	fmt.Fprintf(w, "<h1>Editing: %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>", title, title, string(p.Body))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/delete/"):]
	filePath, err := filepath.Abs(title + ".txt")
	if err != nil {
		return
	}
	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("ERROR: Unable to remove: ", filePath)
	}
}

func main() {
	fmt.Println("Hello World!")

	// Page file
	p1 := &Page{Title: "test1", Body: []byte("This is a sample Page1.")}
	p1.save()
	p2, _ := loadPage("test1")
	fmt.Println(string(p2.Body))

	// Web Page
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/read/", readHandler)
	http.HandleFunc("/update/", updateHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
