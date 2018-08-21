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
	Web struct {
		ClientID                string   `json:"client_id"`
		ProjectID               string   `json:"project_id"`
		AuthURI                 string   `json:"auth_uri"`
		TokenURI                string   `json:"token_uri"`
		AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		RedirectUris            []string `json:"redirect_uris"`
		JavascriptOrigins       []string `json:"javascript_origins"`
	} `json:"web"`
}

// SeverHTTP:Processing of HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
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
		google.New(cs.Web.ClientID, cs.Web.ClientSecret, cs.Web.RedirectUris[0]),
	)

	r := newRoom()
	// show log
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	//chatroom run!
	go r.run()
	// WebServer start
	log.Println("WebServer start. Port:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:err")
	}
}
