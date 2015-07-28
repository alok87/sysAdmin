package controllers

import (
	"os/exec"
	"fmt"
	"io/ioutil"
)

type User struct {
    Username	string	`json:"username"`
    Shelltype	string	`json:"shelltype"`
    HomeFolder	string	`json:"homefolder"`
    Pass		string	`json:"pass"`
    SudoOpt		string	`json:"sudoopt"`
    Operation	string	`json:"operation"`
}

func (user *User) performAction(action string, c *client) {
	var sudoMsg, passMsg string
	var sudoErr, passErr error
	
	//Create or modify user
	createMsg, createErr := user.createOrModifyUser(action)
	if createErr == nil {
		//Update user's passsword
		passMsg, passErr = user.updatePass()
		if passErr != nil {
			passMsg = "Password was not created/updated with the specified password.&#13;&#10;System Message: " + passMsg
		}
		//Make the user a sudo user if option is selected
		sudoErr=nil
		if user.SudoOpt == "Yes" {
			sudoMsg, sudoErr = user.makeSudo()
			if sudoErr != nil {
				sudoMsg = "Unalble to make the user a sudo user.&#13;&#10;System Message: " + sudoMsg
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
}

func (user *User) createOrModifyUser(action string) (string, error) {
	var msgSentBack,actionName string
	
	if action == "Create" {
		actionName = "useradd"
	}else  if action == "Modify" {
		actionName = "usermod"
	}
	
	useraction, lookErr := exec.LookPath(actionName)
    if lookErr != nil {
    	msgSentBack = "Result: Error &#13;&#10;System Message: " + fmt.Sprint(lookErr) + "&#13;&#10;"   
        return msgSentBack, lookErr
    }								    
    cmd := exec.Command("sudo", useraction, "-d", user.HomeFolder, "-s", user.Shelltype, user.Username)
    cmdOut, err := cmd.CombinedOutput()
    if err != nil {
    	msgSentBack = "Result: Validation Error!&#13;&#10;"  + "System Message: " + string(cmdOut) + fmt.Sprint(err) + "&#13;&#10;"   	 
    	return msgSentBack, err
    }
    if string(cmdOut) == "" {
    	msgSentBack = "Result: User, " + user.Username  + " " + action + " Successful !.&#13;&#10;"
    }else {
    	msgSentBack = "Result: User, " + user.Username  + " " + action + " Successful !.&#13;&#10;System Message: " + string(cmdOut) + "&#13;&#10;"
    }
    return msgSentBack, nil
}

func (user *User) deleteUser(c *client) {
	var msgSentBack string
	useraction, lookErr := exec.LookPath("userdel")
    if lookErr != nil {
    	msgSentBack = "Result: Error &#13;&#10;System Message: " + fmt.Sprint(lookErr) + "&#13;&#10;"   
        c.send <- []byte(msgSentBack)
        return
    }	
    cmd := exec.Command("sudo", useraction,	user.Username)	
    cmdOut, err := cmd.CombinedOutput()
    if err != nil {
    	msgSentBack = "Result: Validation Error!&#13;&#10;"  + "System Message: " + string(cmdOut) + fmt.Sprint(err) + "&#13;&#10;"   	 
    	c.send <- []byte(msgSentBack)
    	return
    }		
    if string(cmdOut) == "" {
    	msgSentBack = "Result: User, " + user.Username  + " " +  " Deleted Successfully !.&#13;&#10;"
    }else {
    	msgSentBack = "Result: User, " + user.Username  + " " +  " Deleted Successfully !.&#13;&#10;System Message: " + string(cmdOut) + "&#13;&#10;"
    }
    c.send <- []byte(msgSentBack)    
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
    return msgSentBack, nil
}	