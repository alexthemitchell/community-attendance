package commands

import (
	"database/sql"
	"time"

	"github.com/k0kubun/pp"
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

func (i *importCommand) run(c *kingpin.ParseContext) error {
	db, err := sql.Open("sqlite3", i.dbFileName)
	if err != nil {
		return errors.Wrapf(err, "error opening DB file: %#v", i.dbFileName)
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
	eventTime, err := time.Parse("January _2, 2006 3:04PM PST", i.eventTime)
	if err != nil {
		return errors.Wrap(err, "unable to parse time")
	}

	event := models.NewEvent(i.eventName, "", &eventTime)
	attendance, errs := reader.ReadFile(i.fileName, event)
	err = storage.UpsertEvent(event)
	if err != nil {
		return errors.Wrap(err, "error upserting event")
	}
	if len(errs) > 0 {
		return errors.Wrap(errs[0], "error reading from file")
	}
	for _, entry := range attendance {
		err = storage.UpsertAttendee(entry.Attendee())
		if err != nil {
			log.WithField("entry", entry).WithError(err).Error("error upserting attendee")
			return errors.Wrap(err, "error upserting attendee")
		}
		pp.Println(entry)

	}
	return nil
}

func AddImportSubcommand(app *kingpin.Application) {
	c := app.Command("import", "import attendance information from another source")
	ic := &importCommand{}
	f := c.Command("file", "import information from a local file").Action(ic.run)
	f.Flag("event-time", "the time and date of the event").Required().StringVar(&ic.eventTime)
	f.Arg("file-name", "the name of the file to read").Required().StringVar(&ic.fileName)
	f.Arg("db-file-name", "the name of the sqlite db file").Required().StringVar(&ic.dbFileName)
}
