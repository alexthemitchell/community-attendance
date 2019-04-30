package storage

import (
	"database/sql"
	"fmt"
	gotime "time"

	"github.com/pkg/errors"

	"alexthemitchell/attendance/models"
)

const (
	countEventsQuery           = "SELECT COUNT(*) FROM events"
	createEventsTableStatement = "CREATE TABLE IF NOT EXISTS events (name varchar(255) not null, time DATETIME not null, id int auto_increment primary key, UNIQUE(id))"
	insertEventStatement       = "INSERT INTO events(name, time) VALUES (?,?)"
	deleteEventStatement       = "DELETE FROM events WHERE id=?"
	selectEventStatement       = "SELECT name, id, time FROM events WHERE id=?"
	selectAllEventsStatement   = "SELECT name, id, time FROM events"
	updateEventStatement       = "UPDATE events SET name=?, time=? WHERE id=?"
)

var (
	ErrNoEntryWithEventID = errors.New("no entry exists with the given event identifier")
)

func (s *SQLStorage) GetAllEvents() ([]*models.Event, error) {
	stmt, err := s.db.Prepare(selectAllEventsStatement)
	if err != nil {
		return nil, errors.Wrap(err, "error preparing get all events query")
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.Wrap(err, "error querying for all events")
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		event, err := scanEventFromRow(rows)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning event from row")
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *SQLStorage) CountEvents() (uint, error) {
	stmt, err := s.db.Prepare(countEventsQuery)
	if err != nil {
		return 0, errors.Wrap(err, "error preparing events count query")
	}
	result, err := stmt.Query()
	if err != nil {
		return 0, errors.Wrap(err, "error querying for events count")
	}
	if !result.Next() {
		return 0, errors.New("unexpected SQL result")
	}
	var count uint
	result.Scan(&count)

	return count, nil
}

func (s *SQLStorage) CreateEventsTable() error {
	stmt, err := s.db.Prepare(createEventsTableStatement)
	if err != nil {
		return errors.Wrap(err, "error preparing events table creation")
	}
	_, err = stmt.Exec()
	if err != nil {
		return errors.Wrap(err, "error creating events table")
	}
	return nil
}

func (s *SQLStorage) UpsertEvent(event *models.Event) error {
	if event.ID() == "" {
		if err := s.CreateEvent(event); err != nil {
			return errors.Wrap(err, "error creating event in SQL DB")
		}
	}
	if err := s.UpdateEvent(event); err != nil {
		return errors.Wrap(err, "error upserting event")
	}
	return nil
}

func scanEventFromRow(rows *sql.Rows) (*models.Event, error) {
	var id string
	var name string
	var time string
	rows.Scan(&name, &id, &time)

	fmt.Printf("Read: %#v %#v %#v\n", id, name, time)

	eventTime, err := gotime.Parse(sqlTimestampFormat, time)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing SQL timestamp %#v", time)
	}
	return models.NewEvent(name, id, &eventTime), nil
}

func (s *SQLStorage) FetchEvent(eventID string) (*models.Event, error) {
	stmt, err := s.db.Prepare(selectEventStatement)
	if err != nil {
		return nil, errors.Wrap(err, "error while preparing select statement")
	}
	rows, err := stmt.Query(eventID)
	if err != nil {
		return nil, errors.Wrap(err, "error while executing select statement")

	}
	defer rows.Close()

	exists := rows.Next()
	if !exists {
		return nil, errors.Wrapf(ErrNoEntryWithEventID, "error fetching attendee with ID %#v", eventID)
	}
	return scanEventFromRow(rows)
}

func (s *SQLStorage) CreateEvent(event *models.Event) error {
	stmt, err := s.db.Prepare(insertEventStatement)
	if err != nil {
		return errors.Wrap(err, "error while preparing insert statement")
	}
	eventTime := event.Time().Format(sqlTimestampFormat)
	_, err = stmt.Exec(event.Name(), eventTime)
	if err != nil {
		return errors.Wrap(err, "error while executing insert statement")

	}
	return nil
}

func (s *SQLStorage) UpdateEvent(event *models.Event) error {
	stmt, err := s.db.Prepare(updateEventStatement)
	if err != nil {
		return errors.Wrap(err, "error while preparing update statement")
	}
	eventTime := event.Time().Format(sqlTimestampFormat)
	_, err = stmt.Exec(event.Name(), eventTime, event.ID())
	if err != nil {
		return errors.Wrap(err, "error while executing update statement")

	}
	return nil
}

func (s *SQLStorage) DeleteEvent(eventID string) error {
	stmt, err := s.db.Prepare(deleteEventStatement)
	if err != nil {
		return errors.Wrap(err, "error while preparing delete statement")
	}
	_, err = stmt.Exec(eventID)
	if err != nil {
		return errors.Wrap(err, "error while executing delete statement")

	}
	return nil
}
