package controllers

import (
	"net/http"
	"text/template"
	"sync"
	"path/filepath"
	"strings"
	"os"
	"bufio"
	"log"
	"time"
	"fmt"
	
	"github.com/gorilla/websocket"
)

var (
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second	
)

func Register(template *template.Template) {
	
	uc := new(usersController)
	uc.template = template.Lookup("users.html")
	http.HandleFunc("/users", uc.serveUsers)
	http.HandleFunc("/ws", serveWs)
	
	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("inside ws")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	
	client := &client{
		socket: ws,
		send: make(chan []byte, 256),
		forward: make(chan []byte),
	}
	
	go client.write()
	go client.action()
	client.read()
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
	basePath, err := filepath.Abs("../src/github.com/alok87/sysAdmin/templates")
	if err != nil {
		fmt.Println("Could not find path for templates -", basePath)
		panic(err)
	}
	path := basePath + req.URL.Path
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