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
	"encoding/json"
	
	"github.com/gorilla/websocket"
)

func Register(template *template.Template) {
	
	uc := new(usersController)
	uc.template = template.Lookup("users.html")
	http.HandleFunc("/users", uc.serveUsers)
	http.HandleFunc("/ws", serveWs)
	
	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
}

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

type client struct {
	socket *websocket.Conn
	send chan []byte
	forward chan []byte
}

type User struct {
    Username	string	`json:"username"`
    Shelltype	string	`json:"shelltype"`
    HomeFolder	string	`json:"homefolder"`
    Pass		string	`json:"pass"`
    SudoOpt		string	`json:"sudoopt"`
    Operation	string	`json:"operation"`
}

type Validator struct {
	err error
}

func (v *Validator) MustBeNotEmpty(value string) bool {
	if v.err != nil {
		return true
	}
	if value == "" {
		v.err = fmt.Errorf("Must not be Empty")
		return true
	}
	return false
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside ws")
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
	
	/*
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
	} */
}	

func (c *client) action() {
	var msgSentBack string
	for {
			msg := <- c.forward 
			user := &User{}
			json.Unmarshal([]byte(msg), &user)
			
			//General Form Validation
			validator := new(Validator)
			// Nothing should be empty
			if validator.MustBeNotEmpty(user.Username) || validator.MustBeNotEmpty(user.Operation) || validator.MustBeNotEmpty(user.HomeFolder) || validator.MustBeNotEmpty(user.Pass) || validator.MustBeNotEmpty(user.Shelltype) || validator.MustBeNotEmpty(user.SudoOpt) {
				msgSentBack = "Validation error - all fields required as input. "
				c.send <- []byte(msgSentBack)
				return
				//fmt.Println("from forward > send")
			} else {
				switch user.Operation {
					case "Create": 
									binary, lookErr := exec.LookPath("useradd")
								    if lookErr != nil {
								        panic(lookErr)
								    }
								    cmd := exec.Command("sudo", binary, "-d", user.HomeFolder, "-s", user.Shelltype, user.Username)
								    cmdOut, err := cmd.CombinedOutput()
								    if err != nil {
								    	msgSentBack = "Validation error - user not created. "  + string(cmdOut) + fmt.Sprint(err) 
								    	c.send <- []byte(msgSentBack) 
								    	return
								    }	
								    msgSentBack = "Success - user created. "  + string(cmdOut)	
								    c.send <- []byte(msgSentBack) 				
					case "Modify":
					
					case "Delete":
					
					default:
						msgSentBack = "Error - Operation can be create, delete or modify only. "
						c.send <- []byte(msgSentBack)
				}
			}
	}
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.forward <- msg
			fmt.Println("from socket > forward", msg)
		} else {
			fmt.Println("from socket failed > forward", err)
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println("from send failed > socket", err)
			break
		}
		fmt.Println("from send > socket", msg)
	}
	c.socket.Close()
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