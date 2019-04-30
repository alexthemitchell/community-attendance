package reader

import (
	"encoding/csv"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/alexthemitchell/attendance/models"
)

const rsvpTimeLayout = "January _2, 2006 3:04 PM"
const joinedDateLayout = "January _2, 2006"

func ParseAttendanceFromFile(fileName string, event *models.Event) ([]*models.Attendance, []error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, []error{
			errors.Wrapf(err, "error opening file for read: %#v", fileName),
		}
	}
	defer file.Close()

	reader := csv.NewReader(file)

	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	readData, err := reader.ReadAll()
	if err != nil {
		return nil, []error{
			errors.Wrapf(err, "error reading tab separated data from file: %#v", fileName),
		}
	}

	var attendances []*models.Attendance
	var errs []error
	for rowIndex, row := range readData {
		var preferredName string
		var userID string
		var legalName string
		var profileURL *url.URL
		var isHost bool
		var rsvp bool
		var rsvpTime time.Time
		var joinedDate time.Time
		for cellIndex, cellValue := range row {
			switch cellIndex {
			case 0:
				preferredName = cellValue
				break
			case 1:
				userID = cellValue
				break
			case 3:
				isHost = (cellValue == "Yes")
				break
			case 4:
				rsvp = (cellValue == "Yes")
				break

			case 6:
				rsvpTime, err = time.Parse(rsvpTimeLayout, cellValue)
				if err != nil {
					errs = append(errs, errors.Wrapf(err, "error parsing RSVP time for row %d (%#v)", rowIndex, preferredName))
					continue
				}
				break
			case 7:
				joinedDate, err = time.Parse(joinedDateLayout, cellValue)
				if err != nil {
					errs = append(errs, errors.Wrapf(err, "error parsing RSVP time for row %d (%#v)", rowIndex, preferredName))
					continue
				}
				break
			case 8:
				profileURL, err = url.Parse(cellValue)
				if err != nil {
					errs = append(errs, errors.Wrapf(err, "error parsing profile URL for row %d (%#v)", rowIndex, preferredName))
					continue
				}
				break

			case 9:
				legalName = cellValue
				break
			}
		}
		attendee := models.NewAttendee(preferredName, legalName, userID, profileURL, &joinedDate, isHost)
		att := models.NewAttendance(attendee, event, rsvp, &rsvpTime)
		attendances = append(attendances, att)
	}

	return attendances, errs
}

/*

0 Name
1 User ID
2 Title
3 Event Host
4 RSVP
5 Guests
6 RSVPed on
7 Joined Group on
8 URL of Member Profile
9 (Mandatory) Please provide your name exactly as it appears on your ID here, even if we already have it. Due to security considerations, you cannot be admitted without completing this field.

*/
