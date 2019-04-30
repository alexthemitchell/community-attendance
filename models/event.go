package models

import "time"

type Event struct {
	name string
	id   string
	time *time.Time
}

func NewEvent(name, id string, time *time.Time) *Event {
	return &Event{
		name: name,
		id:   id,
		time: time,
	}
}

func (e *Event) Name() string {
	return e.name
}

func (e *Event) Time() *time.Time {
	return e.time
}

func (e *Event) ID() string {
	return e.id
}
