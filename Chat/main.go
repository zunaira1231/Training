package main

import (
	"Chat/code"
	"flag"
	"github.com/stretchr/objx"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"

)

const (
	Address string = ":3300"
)


// templ represents a single template
type templateHandler struct {
	//for compiling the template once
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates",
		t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,}

	if authCookie, err := r.Cookie("auth"); err == nil {	data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
		t.templ.Execute(w, data)
	}

func main() {



	var addr = flag.String("addr", ":3300", "The addr of the application.")
	flag.Parse() // parse the flags

	// setup the providers
	gomniauth.SetSecurityKey("12345")
	gomniauth.WithProviders(
		github.New("f95a0f0b919ff26e9334", "aaee9a5af90d5125ba15b739e032c8c3805b8933", "http://localhost:3300/auth/callback/github"),
		google.New("572772238310-f90nstftjroa70kh87dhtls99rfjrqic.apps.googleusercontent.com", "4Iad3j2AhrS9oF_LMG5AUQOV", "http://localhost:3300/auth/callback/google"),
		facebook.New("2368189230175016", "71a49ad379b6650b02d43e5ec2f809ad", "http://localhost:3300/auth/callback/facebook"),
	)

	//create new room
	r := code.NewRoom()
	r.Tracer_ = trace.New(os.Stdout)

	http.Handle("/",code.MustAuth( &templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", code.LoginHandler)
	http.Handle("/room", r)

	// get the room going
	go r.Run()

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

