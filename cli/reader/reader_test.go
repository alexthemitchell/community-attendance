package reader

import (
	"testing"
	"time"

	"github.com/alexthemitchell/attendance/models"
	"github.com/stretchr/testify/assert"
)

const (
	happyPathSourceFile = "./test_files/validexample.tsv"
)

func TestParseAttendanceFromFileHappyPath(t *testing.T) {
	now := time.Now()
	event := models.NewEvent("Test Event", "1234", &now)
	attendance, errs := ParseAttendanceFromFile(happyPathSourceFile, event)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, 4, len(attendance))

}
