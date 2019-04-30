package commands

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"alexthemitchell/attendance/models"
	"alexthemitchell/attendance/reader"
	"alexthemitchell/attendance/storage/sql"
)

var log = logrus.StandardLogger()

type importCommand struct {
	fileName   string
	eventName  string
	eventTime  string
	dbFileName string
}

func persistImport(dbName string, event *models.Event, records []*models.Attendance) error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return errors.Wrapf(err, "error opening DB file: %#v", dbName)
	}
	storage, err := storage.NewSQLStorage(db)
	if err != nil {
		return errors.Wrapf(err, "error initializing SQL storage")
	}
	defer storage.Close()
	err = storage.CreateAttendeesTable()
	if err != nil {
		return errors.Wrap(err, "error creating attendees table")
	}
	err = storage.UpsertEvent(event)
	if err != nil {
		return errors.Wrap(err, "error upserting event")
	}
	for _, record := range records {
		err = storage.UpsertAttendee(record.Attendee())
		if err != nil {
			log.WithField("attendee", record.Attendee()).WithError(err).Error("error upserting attendee")
			return errors.Wrap(err, "error upserting attendee")
		}
	}
	return nil
}

func (i *importCommand) run(c *kingpin.ParseContext) error {
	eventTime, err := time.Parse("January _2, 2006 3:04PM PST", i.eventTime)
	if err != nil {
		return errors.Wrap(err, "unable to parse time")
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "unable to create random UUID for event")
	}

	event := models.NewEvent(i.eventName, uuid.String(), &eventTime)
	attendance, errs := reader.ReadFile(i.fileName, event)
	if len(errs) > 0 {
		return errors.Wrap(errs[0], "error reading from file")
	}

	fmt.Printf("processed %d attendance records\n", len(attendance))
	if i.dbFileName != "" {
		persistImport(i.dbFileName, event, attendance)
		fmt.Printf("saved to SQLiteDB %#v\n", i.dbFileName)
	}
	return nil
}

func AddImportSubcommand(app *kingpin.Application) {
	c := app.Command("import", "import attendance information from another source")
	ic := &importCommand{}
	f := c.Command("file", "import information from a local file").Action(ic.run)
	f.Arg("event-name", "the name of the event").Required().StringVar(&ic.fileName)
	f.Arg("event-time", "the time and date of the event").Required().StringVar(&ic.eventTime)
	f.Arg("file-name", "the name of the file to read").Required().StringVar(&ic.fileName)
	f.Flag("save-sqlite", "save the data in a sqlite db file with the provided name").StringVar(&ic.dbFileName)
}
