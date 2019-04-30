package storage

import (
	"database/sql"

	"github.com/pkg/errors"
)

type SQLStorage struct {
	db *sql.DB
}

func (s *SQLStorage) Close() {
	s.db.Close()
}

func NewSQLStorage(db *sql.DB) (*SQLStorage, error) {
	storage := &SQLStorage{db: db}
	return storage, storage.init()
}

func (s *SQLStorage) init() error {
	err := s.CreateAttendeesTable()
	if err != nil {
		return errors.Wrap(err, "error creating attendees table")
	}
	err = s.CreateEventsTable()
	if err != nil {
		return errors.Wrap(err, "error creating events table")

	}
	return nil
}
