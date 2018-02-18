package main

import (
	"io/ioutil"
  "log"
  "net/http"
  "html/template"
)

// a template var to cache the loaded templates to avoid loading them
// on every page view. the .must just makes it freak out and error if
// it cant load the templates, since that would wreck our shit and the
// site wouldnt run, so we want to error
// ALL our templates have be fed to ParseFiles
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))


//Page "object", this is my active record object
type Page struct {
  Title string
  Body []byte
}

//Page object class method to save it to my "database"
func (p *Page) save() error {
    filename := p.Title + ".txt"
    // a standard library function that writes a byte slice to a file
    // third value is the permissions for the file, see Unix docs
    return ioutil.WriteFile(filename, p.Body, 0600)
}

//OPEN SAVED FILE
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

//ACTIONS, controller
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, err := loadPage(title)
    if err != nil {
      http.Redirect(w, r, "/edit/"+title, http.StatusFound)
      return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

//to be DRY, this is the part that loads and executes the files as templates
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

//ROUTES
func main() {
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
