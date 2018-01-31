package main

import (
  "net/http"
  "path"
  "html/template"
  "io/ioutil"

  "github.com/gorilla/mux"
)

// Template Routes
var templates = "templates/general/"
var blogposts = "templates/blogposts/"
var projectposts = "templates/projectposts/"

// Construct Webpage
type Webpage struct {
	Title string
	Body  []byte
}

// Load webpage, else throw an error.
func loadWebpage(htmlfile string, folderpath string) (*Webpage, error){
  filename := folderpath+htmlfile
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Webpage{Title: htmlfile, Body: body}, nil
}

// Render Template
func renderTemplate(w http.ResponseWriter, title string, folderpath string, page *Webpage) {
	temp, _ := template.ParseFiles(folderpath+title)
	temp.Execute(w, page)
}

// Route handler.
// HTML file name must match the route name.
// Routing uses mux.vars to find dynamic route info
func pageHandler(w http.ResponseWriter, r *http.Request) {
  htmlfile := ""
  folderpath := ""
  vars := mux.Vars(r)
  if path.Base(r.URL.Path[:]) == "/" {
    htmlfile = "index.html"
    folderpath = templates
  } else if len(vars) == 0 {
      htmlfile = path.Base(r.URL.Path[:])+".html"
      folderpath = templates
  } else {
      if vars["projectname"] != "" {
        htmlfile = path.Base(r.URL.Path[:])+".html"
        folderpath = projectposts
      } else if vars["postname"] != "" {
        htmlfile = path.Base(r.URL.Path[:])+".html"
        folderpath = blogposts
      }
  }

  page, err := loadWebpage(htmlfile, folderpath)
  if err != nil {
    http.Redirect(w, r, "/404", http.StatusFound)
    return
  }

  renderTemplate(w, htmlfile, folderpath, page)
}
