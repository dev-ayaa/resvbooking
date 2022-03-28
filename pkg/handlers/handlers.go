package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/dev-ayaa/resvbooking/pkg/driver"
	"github.com/dev-ayaa/resvbooking/pkg/forms"
	"github.com/dev-ayaa/resvbooking/pkg/helpers"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"github.com/dev-ayaa/resvbooking/pkg/render"
	"github.com/dev-ayaa/resvbooking/repository"
	"github.com/dev-ayaa/resvbooking/repository/dbRepository"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Repository struct to store the app Config
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepository
}

var Repo *Repository

// NewRepository  create a new repository
func NewRepository(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{App: a,
		DB: dbRepository.NewPostgresRepository(a, db.PSQL)}

}

func NewHandlers(r *Repository) {
	Repo = r
}

// HomePage home page handlers & give the handlers a receiver
func (rp *Repository) HomePage(wr http.ResponseWriter, rq *http.Request) {
	err := render.Template(wr, "home.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

// AboutPage about page  handler
func (rp Repository) AboutPage(wr http.ResponseWriter, rq *http.Request) {
	err := render.Template(wr, "about.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}

}

//ContactPage handler function
func (rp *Repository) ContactPage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "contact.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//JuniorSuitePage  handler function
func (rp *Repository) JuniorSuitePage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "junior.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//PremiumSuitePage handler function
func (rp *Repository) PremiumSuitePage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "premium.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//DeluxeSuitePage handler function
func (rp *Repository) DeluxeSuitePage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "deluxe.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//PenthousePage handler function
func (rp *Repository) PenthousePage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "penthouse.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//ExecutivePage handler function
func (rp *Repository) ExecutivePage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "executive.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//MakeReservationPage handlers function
func (rp *Repository) MakeReservationPage(wr http.ResponseWriter, rq *http.Request) {
	resv, ok := rp.App.Session.Get(rq.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerSideError(wr, errors.New("error linking with sessions"))
		return
	}

	room, err := rp.DB.GetRooms(resv.RoomID)
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}
	resv.Room.RoomName = room.RoomName

	data := make(map[string]interface{})
	stringData := make(map[string]string)

	checkInDate := resv.CheckInDate.Format("2006-01-02")
	checkOutDate := resv.CheckOutDate.Format("2006-01-02")

	stringData["check-in"] = checkInDate
	stringData["check-out"] = checkOutDate

	data["reservation"] = resv

	rp.App.Session.Put(rq.Context(), "reservation", resv)

	render.Template(wr, "make-reservation.page.tmpl", &models.TemplateData{
		Form:       forms.NewForm(nil),
		Data:       data,
		StringData: stringData,
	}, rq)
}

func (rp *Repository) PostMakeReservationPage(wr http.ResponseWriter, rq *http.Request) {
	/*Clients and Server-side Form Validation is process*/
	resv, ok := rp.App.Session.Get(rq.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerSideError(wr, errors.New("error transferring data to post handler"))

	}
	err := rq.ParseForm()
	if err != nil {
		helpers.ServerSideError(wr, err)
	}
	/*

		dateLayout := "2006-01-02"
		checkIn := rq.Form.Get("check-in")
		checkOut := rq.Form.Get("check-out")
		checkInDate, err := time.Parse(dateLayout, checkIn)
		if err != nil {
			helpers.ServerSideError(wr, err)
			return
		}

		checkOutDate, err := time.Parse(dateLayout, checkOut)
		if err != nil {
			helpers.ServerSideError(wr, err)
			return
		}
	*/

	roomID, err := strconv.Atoi(rq.Form.Get("room_id"))
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}
	// var resv models.Reservation
	resv.FirstName = rq.Form.Get("first-name")
	resv.LastName = rq.Form.Get("last-name")
	resv.Email = rq.Form.Get("email")
	resv.PhoneNumber = rq.Form.Get("phone-number")
	// resv.CheckInDate = rq.Form.Get("check-in")
	// resv.CheckOutDate = rq.Form.Get("check-out")
	// 	RoomID:       roomID,

	// 	FirstName:    rq.Form.Get("first-name"),
	// 	LastName:     rq.Form.Get("last-name"),
	// 	Email:        rq.Form.Get("email"),
	// 	PhoneNumber:  rq.Form.Get("phone-number"),
	// 	CheckInDate:  checkInDate,
	// 	CheckOutDate: checkOutDate,
	// 	RoomID:       roomID,
	// }

	form := forms.NewForm(rq.PostForm)

	form.Require("first-name", "last-name", "phone-number", "email")

	form.ValidLenCharacter("first-name", 3, rq)
	form.ValidLenCharacter("last-name", 3, rq)
	form.ValidEmail("email")

	if !form.FormValid() {
		data := make(map[string]interface{})
		data["reservation"] = resv
		err := render.Template(wr, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, rq)
		if err != nil {
			return
		}
		return
	}

	NewResvervationID, err := rp.DB.InsertReservation(resv)
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}
	restriction := models.RoomRestriction{
		ID:            0,
		RoomID:        roomID,
		ReservationID: NewResvervationID,
		RestrictionID: 1,
		CheckInDate:   resv.CheckInDate,
		CheckOutDate:  resv.CheckOutDate,
	}

	err = rp.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}

	rp.App.Session.Put(rq.Context(), "reservation", resv)

	//redirect the data back to avoid submitting the form more than onece
	http.Redirect(wr, rq, "/make-reservation-data", http.StatusSeeOther)
}

func (rp *Repository) MakeReservationSummary(wr http.ResponseWriter, rq *http.Request) {

	resv, ok := rp.App.Session.Get(rq.Context(), "reservation").(models.Reservation)
	if !ok {
		fmt.Println(ok)
		rp.App.Session.Put(rq.Context(), "error", "session has not reservation")
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
		log.Println("Error transferring Data")
		return
	}
	rp.App.Session.Put(rq.Context(), "reservation", resv)

	data := make(map[string]interface{})
	data["reservation"] = resv

	rp.App.Session.Remove(rq.Context(), "reservation")
	err := render.Template(wr, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	}, rq)
	if err != nil {
		return
	}

}

//CheckAvailabilityPage handler Function
func (rp *Repository) CheckAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "check-availability.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//PostCheckAvailabilityPage handler function
func (rp *Repository) PostCheckAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {
	//getting the posted value from the form with respect to the field
	checkIn := rq.Form.Get("check-in")
	checkOut := rq.Form.Get("check-out")

	//Converting the date in string format to time.Time format
	dateLayout := "2006-01-02"
	checkInDate, err := time.Parse(dateLayout, checkIn)
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}
	checkOutDate, err := time.Parse(dateLayout, checkOut)
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}

	rooms, err := rp.DB.SearchForAvailableRoom(checkInDate, checkOutDate)
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}
	for _, room := range rooms {
		rp.App.InfoLog.Println("Rooms Available :: ", room)
	}

	if len(rooms) == 0 {
		rp.App.InfoLog.Println("NO AVAILABLE ROOMS")
		rp.App.Session.Put(rq.Context(), "errors", "No availale rooms")
		http.Redirect(wr, rq, "/check-availability", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms

	//After checking for available room by date and store it in session
	resv := models.Reservation{
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
	}

	rp.App.Session.Put(rq.Context(), "reservation", resv)

	render.Template(wr, "select-available-room.page.tmpl", &models.TemplateData{
		Data: data,
	}, rq)

	//wr.Write([]byte(fmt.Sprintf("Check-in date is %s\nCheck-out date is %s", checkIn, checkOut)))
}

//create a json struct interfaces
type ResponseJSON struct {
	Name    string `json:"name"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// JsonAvailabilityPage  handler Function
func (rp *Repository) JsonAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {

	myResp := ResponseJSON{
		Name:    "Yusuf Akinleye",
		Ok:      true,
		Message: "Available for freelance",
	}

	//Creating a Json file from struct type
	output, err := json.MarshalIndent(myResp, "", "     ")
	//check for errors
	if err != nil {
		//log.Println(err)
		helpers.ServerSideError(wr, err)
	}

	//this type the browser the type of content it is getting
	wr.Header().Set("Content-type", "application/json")
	wr.Write(output)

}

func (rp *Repository) SelectAvailableRoom(wr http.ResponseWriter, rq *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(rq, "id"))
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}
	resv, ok := rp.App.Session.Get(rq.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerSideError(wr, err)
		return
	}
	resv.RoomID = roomID
	rp.App.Session.Put(rq.Context(), "reservation", resv)
	http.Redirect(wr, rq, "/make-reservation", http.StatusSeeOther)
	//render.Template(wr, "make-reservation.page.tmpl", &models.TemplateData{}, rq)
}
