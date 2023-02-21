package timeh

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	ptime "github.com/yaa110/go-persian-calendar"
)

func Init() {
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
}

func DateToJalali(standardDate string, format string) string {
	items := strings.Split(standardDate+" ", " ")
	items[0] += "-0-0"
	items[1] += ":0:0"
	dateItems := strings.Split(items[0], "-")
	timeItems := strings.Split(items[1], ":")
	y, _ := strconv.Atoi(dateItems[0])
	m, _ := strconv.Atoi(dateItems[1])
	d, _ := strconv.Atoi(dateItems[2])
	h, _ := strconv.Atoi(timeItems[0])
	i, _ := strconv.Atoi(timeItems[1])
	s, _ := strconv.Atoi(timeItems[2])
	var t time.Time = time.Date(y, time.Month(m), d, h, i, s, 0, time.UTC)
	pt := ptime.New(t)
	format = strings.ReplaceAll(format, "Y", strconv.Itoa(pt.Year()))
	format = strings.ReplaceAll(format, "m", strconv.Itoa(int(pt.Month())))
	format = strings.ReplaceAll(format, "d", strconv.Itoa(pt.Day()))
	format = strings.ReplaceAll(format, "H", strconv.Itoa(pt.Hour()))
	format = strings.ReplaceAll(format, "i", strconv.Itoa(pt.Minute()))
	format = strings.ReplaceAll(format, "s", strconv.Itoa(pt.Second()))
	return format
}

func JalaliToDate(standardDate string, format string) string {
	items := strings.Split(standardDate, " ")
	dateItems := strings.Split(items[0], "-")
	timeItems := strings.Split(items[1], ":")
	y, _ := strconv.Atoi(dateItems[0])
	m, _ := strconv.Atoi(dateItems[1])
	d, _ := strconv.Atoi(dateItems[2])
	h, _ := strconv.Atoi(timeItems[0])
	i, _ := strconv.Atoi(timeItems[1])
	s, _ := strconv.Atoi(timeItems[2])
	timeZone, _ := time.LoadLocation("Asia/Tehran")
	var pt ptime.Time = ptime.Date(y, ptime.Month(m), d, h, i, s, 0, timeZone)
	t := pt.Time().UTC()

	format = strings.ReplaceAll(format, "Y", strconv.Itoa(t.Year()))
	format = strings.ReplaceAll(format, "m", strconv.Itoa(int(t.Month())))
	format = strings.ReplaceAll(format, "d", strconv.Itoa(t.Day()))
	format = strings.ReplaceAll(format, "H", strconv.Itoa(t.Hour()))
	format = strings.ReplaceAll(format, "i", strconv.Itoa(t.Minute()))
	format = strings.ReplaceAll(format, "s", strconv.Itoa(t.Second()))
	return format
}

func TimeFormat(t time.Time, format string) string {
	Y := fmt.Sprintf("%04d", t.Year())
	format = strings.ReplaceAll(format, "Y", Y)
	m := fmt.Sprintf("%02d", t.Month())
	format = strings.ReplaceAll(format, "m", m)
	d := fmt.Sprintf("%02d", t.Day())
	format = strings.ReplaceAll(format, "d", d)
	h := fmt.Sprintf("%02d", t.Hour())
	format = strings.ReplaceAll(format, "H", h)
	i := fmt.Sprintf("%02d", t.Minute())
	format = strings.ReplaceAll(format, "i", i)
	s := fmt.Sprintf("%02d", t.Second())
	format = strings.ReplaceAll(format, "s", s)

	// TODO : add other formats
	//..........

	return format
}
