package models

import (
	"net/url"
	"time"
)

type Attendee struct {
	preferredName string
	legalName     string
	userID        string
	profileURL    *url.URL
	isHost        bool
	joinedDate    *time.Time
}

func (a *Attendee) PreferredName() string {
	return a.preferredName
}

func (a *Attendee) LegalName() string {
	return a.legalName
}
func (a *Attendee) UserID() string {
	return a.userID
}
func (a *Attendee) ProfileURL() *url.URL {
	return a.profileURL
}
func (a *Attendee) IsHost() bool {
	return a.isHost
}
func (a *Attendee) JoinedDate() *time.Time {
	return a.joinedDate
}

func NewAttendee(preferredName, legalName, userID string, profileURL *url.URL, joinedDate *time.Time, isHost bool) *Attendee {
	return &Attendee{
		preferredName: preferredName,
		legalName:     legalName,
		userID:        userID,
		profileURL:    profileURL,
		joinedDate:    joinedDate,
		isHost:        isHost,
	}
}
