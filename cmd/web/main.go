package main

import (
	"log"
	"net/http"
	"time"

	"github.com/MuiNV/bookings/pkg/config"
	"github.com/MuiNV/bookings/pkg/handlers"
	"github.com/MuiNV/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNo = ":8080"

var session *scs.SessionManager
var app config.AppConfig

func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create templateCache")
	}

	app.TemplateCache = tc
	app.Usecache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	render.NewTemplates(&app)

	srv := &http.Server{
		Addr:    portNo,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
