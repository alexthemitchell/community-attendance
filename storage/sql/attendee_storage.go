package storage

import (
	"database/sql"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/alexthemitchell/community-attendance/models"
)

const (
	sqlTimestampFormat = "2006-01-02T15:04:05Z"

	countAttendeesQuery           = "SELECT COUNT(*) FROM attendees"
	createAttendeesTableStatement = "CREATE TABLE IF NOT EXISTS attendees (preferred_name varchar(255), legal_name varchar(255), user_id varchar(255), profile_url varchar(1000), is_host boolean, joined_date DATETIME, UNIQUE(user_id))"
	insertAttendeeStatement       = "INSERT INTO attendees(preferred_name, legal_name, user_id, profile_url, is_host, joined_date) VALUES (?,?,?,?,?,?)"
	deleteAttendeeStatement       = "DELETE FROM attendees WHERE user_id=?"
	selectAttendeeStatement       = "SELECT * FROM attendees WHERE user_id=?"
	selectAllAttendeesStatement   = "SELECT * FROM attendees"
	updateAttendeeStatement       = "UPDATE attendees SET preferred_name=?, legal_name=?, profile_url=?, is_host=?, joined_date=? WHERE user_id=?"
)

var (
	log                  = logrus.StandardLogger()
	ErrNoEntryWithUserID = errors.New("no entry exists with the given user ID")
)

func (s *SQLStorage) GetAllAttendees() ([]*models.Attendee, error) {
	stmt, err := s.db.Prepare(selectAllAttendeesStatement)
	if err != nil {
		return nil, errors.Wrap(err, "error preparing attendees count query")
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.Wrap(err, "error querying for attendees count")
	}
	var attendees []*models.Attendee
	for rows.Next() {
		attendee, err := scanAttendeeFromRow(rows)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning attendee from row")
		}
		attendees = append(attendees, attendee)

	}
	return attendees, nil
}

func (s *SQLStorage) CountAttendees() (uint, error) {
	stmt, err := s.db.Prepare(countAttendeesQuery)
	if err != nil {
		return 0, errors.Wrap(err, "error preparing attendees count query")
	}
	result, err := stmt.Query()
	if err != nil {
		return 0, errors.Wrap(err, "error querying for attendees count")
	}
	if !result.Next() {
		return 0, errors.New("unexpected SQL result")
	}
	var count uint
	result.Scan(&count)

	return count, nil
}

func (s *SQLStorage) CreateAttendeesTable() error {
	stmt, err := s.db.Prepare(createAttendeesTableStatement)
	if err != nil {
		return errors.Wrap(err, "error preparing attendees table creation")
	}
	_, err = stmt.Exec()
	if err != nil {
		return errors.Wrap(err, "error creating attendees table")
	}
	return nil
}

func (s *SQLStorage) UpsertAttendee(attendee *models.Attendee) error {
	// TODO: Use a transaction to make this more elegant
	if err := s.CreateAttendee(attendee); err == nil {
		// This was a new entry, exit with no error
		return nil
	}
	// If we error on creation, try update in case it already exists
	if err := s.UpdateAttendee(attendee); err != nil {
		return errors.Wrap(err, "error upserting attendee")
	}
	return nil
}

func scanAttendeeFromRow(rows *sql.Rows) (*models.Attendee, error) {
	var preferred_name string
	var legal_name string
	var user_id string
	var profile_url string
	var is_host bool
	var joined_date string
	rows.Scan(&preferred_name, &legal_name, &user_id, &profile_url, &is_host, &joined_date)

	joinDate, err := time.Parse(sqlTimestampFormat, joined_date)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing SQL timestamp %#v", joined_date)
	}
	profileURL, err := url.Parse(profile_url)
	return models.NewAttendee(preferred_name, legal_name, user_id, profileURL, &joinDate, is_host), nil
}

func (s *SQLStorage) FetchAttendee(userID string) (*models.Attendee, error) {
	stmt, err := s.db.Prepare(selectAttendeeStatement)
	if err != nil {
		return nil, errors.Wrap(err, "error while preparing select statement")
	}
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, errors.Wrap(err, "error while executing select statement")

	}
	defer rows.Close()

	exists := rows.Next()
	if !exists {
		return nil, errors.Wrapf(ErrNoEntryWithUserID, "error fetching attendee with ID %#v", userID)
	}
	return scanAttendeeFromRow(rows)
}

func (s *SQLStorage) CreateAttendee(attendee *models.Attendee) error {
	stmt, err := s.db.Prepare(insertAttendeeStatement)
	if err != nil {
		return errors.Wrap(err, "error while preparing insert statement")
	}
	_, err = stmt.Exec(attendee.PreferredName(), attendee.LegalName(), attendee.UserID(), attendee.ProfileURL().String(), attendee.IsHost(), attendee.JoinedDate().Format(sqlTimestampFormat))
	if err != nil {
		return errors.Wrap(err, "error while executing insert statement")

	}
	return nil
}

func (s *SQLStorage) UpdateAttendee(attendee *models.Attendee) error {
	stmt, err := s.db.Prepare(updateAttendeeStatement)
	if err != nil {
		return errors.Wrap(err, "error while preparing update statement")
	}
	joinDate := attendee.JoinedDate().Format(sqlTimestampFormat)
	_, err = stmt.Exec(attendee.PreferredName(), attendee.LegalName(), attendee.ProfileURL().String(), attendee.IsHost(), joinDate, attendee.UserID())
	if err != nil {
		return errors.Wrap(err, "error while executing update statement")

	}
	return nil
}

func (s *SQLStorage) DeleteAttendee(userID string) error {
	stmt, err := s.db.Prepare(deleteAttendeeStatement)
	if err != nil {
		return errors.Wrap(err, "error while preparing delete statement")
	}
	_, err = stmt.Exec(userID)
	if err != nil {
		return errors.Wrap(err, "error while executing delete statement")

	}
	return nil
}
