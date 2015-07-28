package controllers

import (
	"fmt"
	"encoding/json"
	
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send chan []byte
	forward chan []byte
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

func (c *client) action() {
	for {
			msg := <- c.forward 
			user := &User{}
			json.Unmarshal([]byte(msg), &user)
			
			//General Form Validation
			validator := new(Validator)
			
			// Nothing should be empty
			if validator.MustBeNotEmpty(user.Username) || validator.MustBeNotEmpty(user.Operation) || validator.MustBeNotEmpty(user.HomeFolder) || validator.MustBeNotEmpty(user.Pass) || validator.MustBeNotEmpty(user.Shelltype) || validator.MustBeNotEmpty(user.SudoOpt) {
				msgSentBack := "Result: Validation error - All fields required as input.&#13;&#10;"
				//fmt.Println("from forward > send")
				c.send <- []byte(msgSentBack)
			} else {
				switch user.Operation {
					case "Create": 
									user.performAction(user.Operation, c)			 
					case "Modify": 
									user.performAction(user.Operation, c)
					case "Delete":
									user.deleteUser(c)
					default:
						msgSentBack := "Result: Error&#13;&#10;Operation can be - create, delete or modify only.&#13;&#10;"
						c.send <- []byte(msgSentBack)
				}
			}
	}
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.forward <- msg
			//fmt.Println("from socket > forward", msg)
		} else {
			//fmt.Println("from socket failed > forward", err)
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			//fmt.Println("from send failed > socket", err)
			break
		}
		//fmt.Println("from send > socket", msg)
	}
	c.socket.Close()
}