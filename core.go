package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// a template var to cache the loaded templates to avoid loading them
// on every page view. the .must just makes it freak out and error if
// it cant load the templates, since that would wreck our shit and the
// site wouldnt run, so we want to error
// ALL our templates have be fed to ParseFiles
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

//Page "object", this is my active record object
type Page struct {
	gorm.Model
	Title   string
	Body    []byte
	User    *User
	User_id int64
}

type User struct {
	gorm.Model
	Username      string
	password_hash string
	session_token int64
}

//Page object class method to save it to postgres using gorm
func (p *Page) save() error {
	db, err := gorm.Open("postgres", "user=postgres dbname=flexproject sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Create(&Page{Title: p.Title, Body: p.Body})
	return nil
}

//OPEN SAVED FILE
func loadPage(title string) (*Page, error) {
	db, err := gorm.Open("postgres", "user=postgres dbname=flexproject sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	//when you open something in gorm, it is asigned the to variable passed
	//in NOT given as a return value
	var page Page
	db.First(&page, "Title = ?", title)
	return &Page{Title: title, Body: page.Body}, nil
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
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

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
	db, err := gorm.Open("postgres", "user=postgres dbname=flexproject sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Page{})
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
