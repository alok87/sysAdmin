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
	"io/ioutil"
	
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
}	

func (c *client) action() {
	var sudoMsg, passMsg string
	var sudoErr, passErr error
	for {
			msg := <- c.forward 
			user := &User{}
			json.Unmarshal([]byte(msg), &user)
			
			//General Form Validation
			validator := new(Validator)
			// Nothing should be empty
			if validator.MustBeNotEmpty(user.Username) || validator.MustBeNotEmpty(user.Operation) || validator.MustBeNotEmpty(user.HomeFolder) || validator.MustBeNotEmpty(user.Pass) || validator.MustBeNotEmpty(user.Shelltype) || validator.MustBeNotEmpty(user.SudoOpt) {
				msgSentBack := "Result: Validation error - All fields required as input.&#13;&#10;"
				fmt.Println("from forward > send")
				c.send <- []byte(msgSentBack)
			} else {
				switch user.Operation {
					case "Create": 
									//Create user
									createMsg, createErr := user.createUser()
									if createErr == nil {
										
										//Update user's passsword
										passMsg, passErr = user.updatePass()
										if passErr != nil {
											passMsg = "User created, but failed to create password for the new user.&#13;&#10;System Message: " + passMsg
										}
										 
										//Make the user a sudo user if option is selected
										sudoErr=nil
										fmt.Println("user op",user.SudoOpt)
										if user.SudoOpt == "Yes" {
											fmt.Println("user opeation yes")
											sudoMsg, sudoErr = user.makeSudo()
											if sudoErr != nil {
												sudoMsg = "Failed to make the user a sudo user.&#13;&#10;System Message: " + sudoMsg
											}
										}
									    
									    if sudoErr==nil && passErr==nil {
									    	c.send <- []byte(createMsg)
									    }else {
									    	msgSentBack := createMsg + "&#13;&#10;" + passMsg + "&#13;&#10;" + sudoMsg
									    	c.send <- []byte(msgSentBack)
									    }
									    
									}else {
										c.send <- []byte(createMsg)
									}			 
					case "Modify":
									
					
					case "Delete":
					
					default:
						msgSentBack := "Result: Error&#13;&#10;Operation can be - create, delete or modify only.&#13;&#10;"
						c.send <- []byte(msgSentBack)
				}
			}
	}
}

func (user *User) createUser() (string, error) {
	var msgSentBack string
	useradd, lookErr := exec.LookPath("useradd")
    if lookErr != nil {
    	msgSentBack = "Result: Error &#13;&#10;System Message: " + fmt.Sprint(lookErr) + "&#13;&#10;"   
        return msgSentBack, lookErr
    }								    
    cmd := exec.Command("sudo", useradd, "-d", user.HomeFolder, "-s", user.Shelltype, user.Username)
    cmdOut, err := cmd.CombinedOutput()
    if err != nil {
    	msgSentBack = "Result: Validation Error!&#13;&#10;"  + "System Message: " + string(cmdOut) + fmt.Sprint(err) + "&#13;&#10;"   	 
    	return msgSentBack, err
    }
    msgSentBack = "User " + user.Username  + " created.&#13;&#10;System Message: " + string(cmdOut) + "&#13;&#10;"
    return msgSentBack, nil
}

func (u *User) updatePass() (string, error) {
	chpasswd, lookErr := exec.LookPath("chpasswd")
	if lookErr != nil {
	     return "Failed to find chpasswd in Path", lookErr
	}
	userPass := u.Username + ":" + u.Pass
	passCmd := exec.Command("sudo", chpasswd)	
	passCmdIn, err := passCmd.StdinPipe()
	if err == nil {
		passCmdOut, outerr := passCmd.StdoutPipe()
			if outerr == nil {
				passCmd.Start()
				passCmdIn.Write([]byte(userPass))
				passCmdIn.Close()
				_, readerr := ioutil.ReadAll(passCmdOut)
				if readerr == nil {
					passCmd.Wait()
					return "Password updated", nil
				} else {
					return "Error reading passCmdOut", readerr
				} 	
			} else {
				return "Error reading STDOUT", outerr
			}
	} else	{
		return "Error reading STDIN", err
	}
}		


func (user *User) makeSudo() (string, error) {
	fmt.Println("Inside makesudo()")
	var msgSentBack string
	usermod, lookErr := exec.LookPath("usermod")
    if lookErr != nil {
    	fmt.Println("Inside llookErr()")
    	msgSentBack = "System Message: " + fmt.Sprint(lookErr) + "&#13;&#10;"   
    	fmt.Println(msgSentBack)
        return msgSentBack, lookErr
    }								    
    cmd := exec.Command("sudo", usermod, "-a", "-G", "wheel", user.Username)
    cmdOut, err := cmd.CombinedOutput()
    if err != nil {
    	msgSentBack = "System Message: " + string(cmdOut) + fmt.Sprint(err) + "&#13;&#10;"   	 
    	fmt.Println(msgSentBack)
    	return msgSentBack, err
    }
    msgSentBack = "User " + user.Username  + " made a sudo user.&#13;&#10;System Message: " + string(cmdOut) + "&#13;&#10;"
    fmt.Println(msgSentBack)
    return msgSentBack, nil
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