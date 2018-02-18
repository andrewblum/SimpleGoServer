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
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
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

//higher order function! we pass the vanilla handler function into this
// and then we return a http.HandlerFunc that calls the original
// vanilla version with the title string validated and passed as an argument 
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

//ROUTES
func main() {
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
