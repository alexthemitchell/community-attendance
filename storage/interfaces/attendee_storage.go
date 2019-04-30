package storage

import (
	"github.com/alexthemitchell/community-attendance/models"
)

type AttendeeStorage interface {
	CountAttendees() (uint, error)
	FetchAttendee(userID string) (*models.Attendee, error)
	UpsertAttendee(attendee *models.Attendee) error
	DeleteAttendee(userID string) error
}
