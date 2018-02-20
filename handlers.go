package main

import (
  "net/http"
  "path"
  "html/template"
  "io/ioutil"
  "path/filepath"
  "os"
  "log"

  "github.com/gorilla/mux"
)

// Template Routes
var templates = "templates/general/"
var blogposts = "templates/blogposts/"
var projectposts = "templates/projectposts/"

// Globals
var htmlfiles, _ = grabHTML()
var parsedtemplates = ParseTemplates(htmlfiles)

type Webpage struct {
	Title string
	Body  []byte
}

// Get a list of HTML files
func grabHTML() ([]string, error) {
    search_dir := []string{templates, projectposts, blogposts}
    files := make([]string, 0)
    for _, item := range search_dir{
      err := filepath.Walk(item, func(path string, f os.FileInfo, err error) error {
              if filepath.Ext(path) == ".html" {
                files = append(files, path)
              }
              return nil
          })
      if err != nil {
          log.Println(err)
      }
  }
  return files, nil
}

// Create cached templates
func ParseTemplates([]string) (*template.Template) {
  temp := template.New("")
  for _, file := range htmlfiles {
		_, err := temp.ParseFiles(file)
    if err != nil {
      log.Println(err)
    }
	}
  return temp
}

// Load webpage, else throw an error.
func loadWebpage(htmlfile string, folderpath string) (*Webpage, error){
  filename := folderpath+htmlfile
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Println(err)
    return nil, err
  }
  return &Webpage{Title: htmlfile, Body: body}, nil
}

// Render Template
func renderTemplate(w http.ResponseWriter, title string, folderpath string, page *Webpage) {
  err := parsedtemplates.ExecuteTemplate(w, title, page)
  if err != nil {
    log.Println(err)
  }
}

// 404
func handle404(w http.ResponseWriter, r *http.Request) {
  htmlfile := "404.html"
  folderpath := templates
  page, err := loadWebpage(htmlfile, folderpath)
  if err != nil {
    log.Println(err)
    http.Redirect(w, r, "/404", http.StatusFound)
    return
  }
  renderTemplate(w, htmlfile, folderpath, page)
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
    log.Println(err)
    http.Redirect(w, r, "/404", http.StatusFound)
    return
  }
  renderTemplate(w, htmlfile, folderpath, page)
}
