package models

import "time"

type Attendance struct {
	attendee *Attendee
	event    *Event
	rsvpTime *time.Time
	rsvp     bool
}

func NewAttendance(attendee *Attendee, event *Event, rsvp bool, rsvpTime *time.Time) *Attendance {
	return &Attendance{
		attendee: attendee,
		event:    event,
		rsvpTime: rsvpTime,
		rsvp:     rsvp,
	}
}

func (a *Attendance) Attendee() *Attendee {
	return a.attendee
}

func (a *Attendance) Event() *Event {
	return a.event
}

func (a *Attendance) RSVPTime() *time.Time {
	return a.rsvpTime
}

func (a *Attendance) RSVP() bool {
	return a.rsvp
}
