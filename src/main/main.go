package main

import (
	"net/http"
	"log"
	"os"
	"text/template"
	
	"controllers"
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
	
	basePath := "templates"
	templateFolder, _ := os.Open(basePath)
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