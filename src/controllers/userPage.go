package controllers

import (
	"net/http"
	"text/template"	
)

type usersController struct {
	template *template.Template
}

func (this *usersController) serveUsers(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/users" {
		http.Error(w, "Not found", 404)
		return
	}
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	this.template.Execute(w, req)
	
}	
	
	

	