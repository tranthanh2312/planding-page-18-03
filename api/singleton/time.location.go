package singleton

import (
	"log"
	"time"
)

var loc *time.Location

func InitTimeLocation() {
	if err := setTimezone("Asia/Bangkok"); err != nil {
		log.Fatal(err) // most likely timezone not loaded in Docker OS
	}

	log.Println("InitTime")
}

func setTimezone(tz string) error {
	location, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}

	loc = location
	return nil
}

func GetTime(t time.Time) time.Time {
	return t.In(loc)
}

func GetLoc() *time.Location {
	return loc
}
