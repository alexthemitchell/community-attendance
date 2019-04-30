package commands

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alexthemitchell/community-attendance/models"
	"github.com/alexthemitchell/community-attendance/storage/sql"
)

const joinDateDisplayFormat = "2006-01-02"

type listAttendeesCommand struct {
	dbFileName string
}

func attendeeLineFormatWithMaxLengths(maxPreferredName, maxLegalName, maxJoinedDate int) string {
	return fmt.Sprintf("%%-%ds\t%%%ds\t%%%ds\t%%s\n", maxPreferredName, maxLegalName, maxJoinedDate)
}

func lineFormatForAttendees(attendees []*models.Attendee) string {
	var maxPreferredNameLength int
	var maxLegalNameLength int
	var maxJoinedDateLength int
	for _, attendee := range attendees {
		maxPreferredNameLength = max(maxPreferredNameLength, len(attendee.PreferredName()))
		maxLegalNameLength = max(maxLegalNameLength, len(attendee.LegalName()))
		maxJoinedDateLength = max(maxJoinedDateLength, len(attendee.JoinedDate().Format(joinDateDisplayFormat)))
	}

	return attendeeLineFormatWithMaxLengths(
		maxPreferredNameLength,
		maxLegalNameLength,
		maxJoinedDateLength)
}

func printAttendeesToScreen(attendees []*models.Attendee) {
	lineFormat := lineFormatForAttendees(attendees)
	fmt.Fprintf(os.Stdout, lineFormat,
		aurora.Bold("Preferred Name"),
		aurora.Bold("Legal Name"),
		aurora.Bold("Joined Date"),
		aurora.Bold("Host?"),
	)
	for _, attendee := range attendees {
		var hostMarker string
		if attendee.IsHost() {
			hostMarker = "Host"
		}
		fmt.Fprintf(os.Stdout, lineFormat,
			attendee.PreferredName(),
			attendee.LegalName(),
			attendee.JoinedDate().Format(joinDateDisplayFormat),
			hostMarker)
	}
}

func (l *listAttendeesCommand) run(c *kingpin.ParseContext) error {
	db, err := sql.Open("sqlite3", l.dbFileName)
	if err != nil {
		return errors.Wrapf(err, "error opening DB file%#v", l.dbFileName)
	}
	storage, err := storage.NewSQLStorage(db)
	if err != nil {
		return errors.Wrapf(err, "error initializing SQL storage")
	}
	defer storage.Close()
	attendees, err := storage.GetAllAttendees()
	if err != nil {
		return errors.Wrap(err, "error getting attendees from storage")
	}
	printAttendeesToScreen(attendees)
	return nil
}
