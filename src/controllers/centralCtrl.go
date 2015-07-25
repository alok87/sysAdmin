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
	"os/exec"
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
	http.HandleFunc("/", uc.serveUsers)
	http.HandleFunc("/ws", serveWs)
	
	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	
	username := r.FormValue("username")
	shelltype := r.FormValue("shelltype")
	homefolder := r.FormValue("homefolder")
	pass := r.FormValue("pass")
	sudoopt := r.FormValue("sudoopt")
	operation := r.FormValue("operation")
	
	binary, lookErr := exec.LookPath("echo")
    if lookErr != nil {
        panic(lookErr)
    }
    cmd := exec.Command(binary, username, shelltype, homefolder, pass, sudoopt, operation)
    cmdOut, err := cmd.Output()
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(cmdOut))  
    //ws.SetWriteDeadline(time.Now().Add(writeWait))
	if err := ws.WriteMessage(websocket.TextMessage, cmdOut); err != nil {
			fmt.Println(err)
			return
	} 
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