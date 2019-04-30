package main

import (
	"fmt"
	"os"
	"time"

	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"

	"alexthemitchell/attendance/reader"
)

var log = logrus.StandardLogger()

const fileName = "Community_Hack_Night.xls"

func main() {
	now := time.Now()
	attendances, errs := reader.ReadFile(fileName, "Community Hack Night May", &now)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	for _, attendance := range attendances {
		pp.Println(attendance.Attendee())
		pp.Println(attendance.RSVPTime())
	}
}
