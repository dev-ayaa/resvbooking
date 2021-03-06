package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"html/template"
	"log"
)

//Avoiding creating templates cache all the time a page is display making sure
//it doesn't import anything but can be access in any part of the application
//use in the render & handlers

type AppConfig struct {
	UseCache     bool
	TempCache    map[string]*template.Template
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
	InProduction bool
	Session      *scs.SessionManager
	MailChannel  chan models.MailData
}
