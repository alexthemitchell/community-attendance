package commands

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alexthemitchell/attendance/models"
	"github.com/alexthemitchell/attendance/storage/sql"
)

const eventTimeDisplayFormat = "2006-01-02 03:04 PM PST"

type listEventsCommand struct {
	dbFileName string
}

func lineFormatForEvents(events []*models.Event) string {
	var maxNameLength int
	var maxTimeLength int
	var maxIDLength int
	for _, event := range events {
		maxNameLength = max(maxNameLength, len(event.Name()))
		maxIDLength = max(maxIDLength, len(event.ID()))
		maxTimeLength = max(maxTimeLength, len(event.Time().Format(eventTimeDisplayFormat)))
	}

	return eventLineFormatWithMaxLengths(
		maxNameLength,
		maxIDLength,
		maxTimeLength)
}
func eventLineFormatWithMaxLengths(maxNameLength, maxIDLength, maxTimeLength int) string {
	return fmt.Sprintf("%%-%ds\t%%%ds\t%%%ds\n", maxNameLength, maxIDLength, maxTimeLength)
}

func printEventsToScreen(events []*models.Event) {
	lineFormat := lineFormatForEvents(events)
	fmt.Fprintf(os.Stdout, lineFormat,
		aurora.Bold("Event"),
		aurora.Bold("Time"),
		aurora.Bold("ID"),
	)
	for _, event := range events {
		fmt.Fprintf(os.Stdout, lineFormat,
			event.Name(),
			event.Time().String(),
			event.ID())
	}
}

func (l *listEventsCommand) run(c *kingpin.ParseContext) error {
	db, err := sql.Open("sqlite3", l.dbFileName)
	if err != nil {
		return errors.Wrapf(err, "error opening DB file%#v", l.dbFileName)
	}
	storage, err := storage.NewSQLStorage(db)
	if err != nil {
		return errors.Wrapf(err, "error initializing SQL storage")
	}
	defer storage.Close()
	events, err := storage.GetAllEvents()
	if err != nil {
		return errors.Wrap(err, "error getting events from storage")
	}
	fmt.Println(len(events))
	printEventsToScreen(events)
	return nil
}
