package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth/providers/github"

	"github.com/stretchr/gomniauth/providers/facebook"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

// template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

type clientSecret struct {
	Google struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectUris []string `json:"redirect_uris"`
	} `json:"google"`
	Facebook struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectUris []string `json:"redirect_uris"`
	} `json:"facebook"`
	Github struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectUris []string `json:"redirect_uris"`
	} `json:"github"`
}

// SeverHTTP:Processing of HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "Application Address")
	flag.Parse() // Interpret the flag

	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalln("jsonファイルの読み込みに失敗しました", err)
	}
	// JSONデコード
	cs := new(clientSecret)
	if err := json.Unmarshal(bytes, &cs); err != nil {
		log.Fatal(err)
	}

	//Gomniauthの設定
	gomniauth.SetSecurityKey("SecurityKey")
	gomniauth.WithProviders(
		google.New(cs.Google.ClientID, cs.Google.ClientSecret, cs.Google.RedirectUris[0]),
		facebook.New(cs.Facebook.ClientID, cs.Facebook.ClientSecret, cs.Facebook.RedirectUris[0]),
		github.New(cs.Github.ClientID, cs.Github.ClientSecret, cs.Github.RedirectUris[0]),
	)

	r := newRoom(UseFileSystemAvatar)
	// show log
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.Handle("/avatars/",
		http.StripPrefix("/avatars",
			http.FileServer(http.Dir("./avatars"))))
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/room", r)

	//chatroom run!
	go r.run()
	// WebServer start
	log.Println("WebServer start. Port:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:err")
	}
}
