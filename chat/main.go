package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// SeverHTTP:Processing of HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	//chatroom run!
	go r.run()
	// WebServer start
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:err")
	}
}
