package storage

import (
	"github.com/alexthemitchell/attendance/models"
)

type AttendeeStorage interface {
	CountAttendees() (uint, error)
	FetchAttendee(userID string) (*models.Attendee, error)
	UpsertAttendee(attendee *models.Attendee) error
	DeleteAttendee(userID string) error
}
