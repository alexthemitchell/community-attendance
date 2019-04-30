package main

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alexthemitchell/community-attendance/cli/commands"
)

func main() {
	app := kingpin.New("attendance", "Event attendee forecasting software")
	commands.AddImportSubcommand(app)
	commands.AddListSubcommand(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
