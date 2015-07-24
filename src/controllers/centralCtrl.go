package controllers

import (
	"net/http"
	"text/template"
	"sync"
	"path/filepath"
	"strings"
	"os"
	"bufio"
)

func Register(template *template.Template) {
	
	uc := new(usersController)
	uc.template = template.Lookup("users.html")
	http.HandleFunc("/", uc.serveUsers)
	
	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
}

type templateHandler struct {
	fileName string
	once sync.Once
	template *template.Template
}

func (this *templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	this.once.Do(func() {
			this.template = template.Must(template.ParseFiles(filepath.Join("templates/", this.fileName)))
		})
	this.template.Execute(w, nil)
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "templates" + req.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	} else {
		contentType = "text/plain"
	}
	
	f, err := os.Open(path)
	
	if err == nil {
		defer f.Close()
		w.Header().Add("Content Type", contentType)
		
		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}