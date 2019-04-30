package commands

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func AddListSubcommand(app *kingpin.Application) {
	c := app.Command("list", "list information from centralized storage")

	lac := &listAttendeesCommand{}
	a := c.Command("attendees", "show list of attendees").Action(lac.run)
	a.Arg("db-file-name", "the name of the sqlite db file").Required().StringVar(&lac.dbFileName)

	lec := &listEventsCommand{}
	e := c.Command("events", "show list of events").Action(lec.run)
	e.Arg("db-file-name", "the name of the sqlite db file").Required().StringVar(&lec.dbFileName)
}
