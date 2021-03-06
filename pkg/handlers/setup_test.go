package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"

	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/dev-ayaa/resvbooking/pkg/helpers"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"github.com/dev-ayaa/resvbooking/pkg/render"
)

var session *scs.SessionManager
var app config.AppConfig

var functions = template.FuncMap{
	//format a dates, currents date

	"dateFormat": render.RenderDateFormat,
	"format":     render.RenderFormat,
	"iterate":    render.RenderIterate,
	"add":        render.RenderAddUp,
}

var templatesPath = "./../../templates"

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.User{})
	gob.Register(models.MailData{})

	app.InProduction = false

	infoLogger := log.New(os.Stdout, "INFO ::\t", log.LstdFlags)
	app.InfoLog = infoLogger

	errorLogger := log.New(os.Stdout, "ERROR ::\t", log.LstdFlags|log.Lshortfile)
	app.ErrorLog = errorLogger

	session = scs.New()
	session.Lifetime = 24 * time.Hour              // how to keep the session of users
	session.Cookie.Persist = true                  //To keep cookies
	session.Cookie.SameSite = http.SameSiteLaxMode //if the user visit the same sites again
	session.Cookie.Secure = app.InProduction       // is the application in production or development

	//Getting the templates cache
	tc, err := TemplateTestCache()
	fmt.Println(tc, err)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Cannot create template cache")
	}

	// storing the cache in the app config
	app.TempCache = tc
	app.Session = session

	//authorize using cache
	app.UseCache = true

	//Referencing the map store in the app AppConfig
	repo := NewTestRepository(&app)
	NewHandlers(repo)
	helpers.NewHelper(&app)

	render.NewTemplates(&app)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	//	mux.Use(NoSurf)

	mux.Use(SessionLoad)

	mux.Get("/", Repo.HomePage)
	mux.Get("/about", Repo.AboutPage)
	mux.Get("/contact", Repo.ContactPage)
	mux.Get("/junior-suite", Repo.JuniorSuitePage)
	mux.Get("/premium-suite", Repo.PremiumSuitePage)
	mux.Get("/deluxe-suite", Repo.DeluxeSuitePage)
	mux.Get("/penthouse-suite", Repo.PenthousePage)
	mux.Get("/executive-suite", Repo.ExecutivePage)

	mux.Get("/make-reservation", Repo.MakeReservationPage)
	mux.Post("/make-reservation", Repo.PostMakeReservationPage)
	mux.Get("/make-reservation-data", Repo.MakeReservationSummary)

	//mux.Get("/check-availability", Repo.CheckAvailabilityPage)
	mux.Get("/check-availability", Repo.CheckAvailabilityPage)
	mux.Post("/check-availability", Repo.PostCheckAvailabilityPage)

	mux.Get("/json-availability", Repo.JsonAvailabilityPage)
	mux.Post("/json-availability", Repo.JsonAvailabilityPage)

	mux.Get("/login", Repo.LoginPage)
	mux.Post("/login", Repo.PostLoginPage)
	mux.Get("/logout", Repo.LogOutPage)

	//setting up the admin page

	//mux.Use(Authenticate)
	mux.Get("/admin/dashboard", Repo.AdminPage)
	mux.Get("/admin/admin-new-reservation", Repo.AdminNewReservation)
	mux.Get("/admin/admin-all-reservation", Repo.AdminAllReservation)
	mux.Get("/admin/admin-reservation-calendar", Repo.AdminReservationCalendar)
	mux.Post("/admin/admin-reservation-calendar", Repo.PostAdminReservationCalendar)

	mux.Get("/admin/admin-show-reservation/{src}/{id}/show", Repo.AdminShowReservation)
	mux.Post("/admin/admin-show-reservation/{src}/{id}", Repo.PostAdminShowReservation)

	mux.Get("/admin/admin-delete-reservation/{src}/{id}/done", Repo.AdminDeleteReservation)
	mux.Get("/admin/admin-process-reservation/{src}/{id}/done", Repo.AdminProcessReservation)

	//This allows files static files like images and icon to display in the html
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

func NoSurf(next http.Handler) http.Handler {
	//Cross sites examined
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{

		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad Loads and save the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// TemplateCache Working with layout and building a template cache
func TemplateTestCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", templatesPath))
	if err != nil {
		return cache, err
	}

	for _, pg := range pages {
		filename := filepath.Base(pg)
		tmp, err := template.New(filename).Funcs(functions).ParseFiles(pg)

		if err != nil {
			return cache, err
		}
		matchTp, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", templatesPath))
		if len(matchTp) > 0 {
			tmp, err = tmp.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", templatesPath))
			if err != nil {
				return cache, err
			}
		}
		cache[filename] = tmp
	}
	return cache, nil

}
