// Simple Golang Website
// by @alfonsodev 17.10.15
// Features:
//  - Template engine examples with pongo2
//  - Routing
//  - Sessions
//  - Middlewares
//  - Google Sing ing and API calls
//  - xsrftoken
// Todo:
// - Randomize cookie store secret and xsfr token init
// - Recover from panic and show error page
// - 401 unAuth page
// - Tests
// - Design

package main

import (
	"code.google.com/p/google-api-go-client/plus/v1"
	"code.google.com/p/xsrftoken"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
)

type Config struct {
	SessionName string
	Port        string
}

var cfg Config          // App configuration
var conf *oauth2.Config // Google configuration
var store = sessions.NewCookieStore([]byte("1--!!@323kjlkb1#@$V3k1jb31S}{23jcl2"))

func readEnvironment(cfg *Config) {
	cfg.SessionName = os.Getenv("SGW_SESSION")
	cfg.Port = os.Getenv("SGW_PORT")

	if cfg.Port == "" {
		cfg.Port = "3030"
	}
	if cfg.SessionName == "" {
		cfg.SessionName = "sgw"
	}

	if os.Getenv("GOOGLE_CLIENT_ID") == "" {
		panic(" GOOGLE_CLIENT_ID is mandatory")
	}

	if os.Getenv("GOOGLE_CLIENT_ID") == "" {
		panic(" GOOGLE_CLIENT_ID is mandatory")
	}

	if os.Getenv("GOOGLE_CLIENT_SECRET") == "" {
		panic(" GOOGLE_CLIENT_SECRET is mandatory")
	}

	if os.Getenv("GOOGLE_CLIENT_REDIRECT") == "" {
		panic(" GOOGLE_CLIENT_REDIRECT is mandatory")
	}
}

func xsrfMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, _ := store.Get(r, cfg.SessionName)
	//TODO aabb234 is not a good value for the function Generate
	session.Values["xsrf"] = xsrftoken.Generate("aabb234", "1", "login")
	next(rw, r)
}

func sessionMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, _ := store.Get(r, cfg.SessionName)
	next(rw, r)
}

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, cfg.SessionName)
	conf = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_CLIENT_REDIRECT"),
		Scopes: []string{
			"https://www.googleapis.com/auth/plus.login",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	if session.Values["googleId"] != nil {
		http.Redirect(rw, r, "/dashboard", 301)
	}
	// Generate google signin url with xsrf token
	url := conf.AuthCodeURL(session.Values["xsrf"].(string))
	tmpl := pongo2.Must(pongo2.FromFile("./templates/home.html"))
	err := tmpl.ExecuteWriter(pongo2.Context{"GoogleAuthUrl": url}, rw)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func CallbackHandler(rw http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, cfg.SessionName)
	values := r.URL.Query()
	code := values["code"]
	tok, err := conf.Exchange(oauth2.NoContext, code[0])
	if err != nil {
		log.Fatal(err)
	}
	client := conf.Client(oauth2.NoContext, tok)
	service, err := plus.New(client)
	if err != nil {
		fmt.Printf("Error: %+v", err)
	}

	people := service.People.Get("me")
	person, err := people.Do()
	if err != nil {
		fmt.Printf("Error: %+v", err)
	}

	session.Values["googleId"] = person.Id
	session.Save(r, rw)

	http.Redirect(rw, r, "/dashboard", 301)
}

func DashboardHandler(rw http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, cfg.SessionName)
	tmpl := pongo2.Must(pongo2.FromFile("./templates/dash.html"))
	if session.Values["googleId"] != "" {
		err := tmpl.ExecuteWriter(pongo2.Context{"GoogleId": session.Values["googleId"]}, rw)
		if err != nil {
			panic(err)
		}
	}

}

func main() {
	readEnvironment(&cfg)
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/callback", CallbackHandler)

	router.HandleFunc("/dashboard", DashboardHandler)
	// router.HandleFunc("/401", UnauthHandler)
	// router.HandleFunc("/settings", SettingsHandler)
	//router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	// router.HandleFunc("/{username}", DashboardHandler)

	n := negroni.New(negroni.NewStatic(http.Dir("public")))
	n.Use(negroni.HandlerFunc(xsrfMiddleware))
	n.Use(negroni.HandlerFunc(sessionMiddleware))
	// n.Use(negroni.HandlerFunc(resolv.usernameMiddleware))
	n.UseHandler(router)
	n.Run(":" + cfg.Port)

	fmt.Printf("\n Listening on port %v ...\n", cfg.Port)
}
