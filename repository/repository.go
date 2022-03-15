package repository

import "github.com/dev-ayaa/resvbooking/pkg/models"

type DatabaseRepository interface {
	AllUser() bool

	InsertReservation(resv models.Reservation) error
}
