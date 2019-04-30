package commands

import (
	"database/sql"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"

	"alexthemitchell/attendance/storage/sql"
)

type exportCommand struct {
	eventName  string
	eventID    string
	dbFileName string
}

func (i *exportCommand) run(c *kingpin.ParseContext) error {
	db, err := sql.Open("sqlite3", i.dbFileName)
	if err != nil {
		return errors.Wrapf(err, "error opening DB file: %#v", i.dbFileName)
	}
	storage, err := storage.NewSQLStorage(db)
	if err != nil {
		return errors.Wrapf(err, "error initializing SQL storage")
	}
	defer storage.Close()
	events, err := storage.GetAllEvents()
	if err != nil {
		return errors.Wrap(err, "error fetching events from storage")
	}
	pp.Println(events)
	return nil
}

func AddExportSubcommand(app *kingpin.Application) {
	c := app.Command("export", "format data for transit to other software")
	ec := &exportCommand{}
	f := c.Command("guestlist", "prepare a single day's event for entry into Envoy").Action(ec.run)
	f.Flag("event-id", "the id of the event as shown on the list command").Required().StringVar(&ec.eventID)
	f.Arg("db-file-name", "the name of the sqlite db file").Required().StringVar(&ec.dbFileName)
}
