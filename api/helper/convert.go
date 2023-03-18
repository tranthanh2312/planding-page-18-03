package helper

import (
	"log"
	"strings"
	"time"

	"github.com/Tech-by-GL/dashboard/singleton"
)

func ConvertMysqlTimeUnixTime(mysqlTime string) int64 {
	res1 := strings.Replace(mysqlTime, "T", " ", 1)
	res2 := res1[:19]

	// YYYY-MM-DD
	layout := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(layout, res2, singleton.GetLoc())
	if err != nil {
		log.Println("date err: ", err)
		panic(err)
	}

	if t.Before(time.Date(2022, 11, 17, 0, 0, 0, 0, singleton.GetLoc())) {
		return t.Unix() + 7*60*60
	}

	return t.Unix()
}

func ConvertUnixTimeMySqlTime(t int64) string {
	tm := time.Unix(t, 0)

	return tm.Format("2006-01-02 15:04:05")
}
