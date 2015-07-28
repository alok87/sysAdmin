package main

import (
	"net/http"
	"log"
	"os"
	"text/template"
	"fmt"
	"path/filepath"
	
	"github.com/alok87/sysAdmin/src/controllers"
)

func main() {	
	templates := populateTemplates()
	controllers.Register(templates)
	
	if err := http.ListenAndServe(":8005", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


func populateTemplates() *template.Template {
	result := template.New("templates")
	
	// basePath := "templates"
	basePath, err := filepath.Abs("../src/github.com/alok87/sysAdmin/templates")
	if err != nil {
		fmt.Println("Could not find path for templates -", basePath)
		panic(err)
	}
	templateFolder, err := os.Open(basePath)
	if err != nil {
		fmt.Println("Could not open templates - ", basePath)
		panic(err)
	}
	defer templateFolder.Close()
	
	templatePathsRaw, _ := templateFolder.Readdir(-1)
	
	templatePaths := new([]string)
	for _, pathInfo := range templatePathsRaw {
		if !pathInfo.IsDir() {
			*templatePaths = append(*templatePaths,
				basePath + "/" + pathInfo.Name())
		}
	}
	
	result.ParseFiles(*templatePaths...)
	
	return result
}