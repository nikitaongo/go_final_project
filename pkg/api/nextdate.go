package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	maxDays int = 400
	layout      = "20060102"
)

// nextDayHandler handles Get-request with repeat pattern and responses next date
func nextDayHandler(res http.ResponseWriter, req *http.Request) {
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	now := req.FormValue("now")
	if now == "" {
		now = time.Now().Format(layout)
	}
	nowTime, err := time.Parse(layout, now)
	if err != nil {
		http.Error(res, fmt.Sprintf("can't parse <now> parameter: %v", err), http.StatusInternalServerError)
		return
	}

	next, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(next))
}

// daysIn calculates the number of days in the passed month of passed year
func daysIn(month, year int) int {
	t := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	return t.Day()
}

// NextDate calculates next date with repeat pattern, now and start dates
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	nextTime, err := time.Parse(layout, dstart)
	if err != nil {
		return "", fmt.Errorf("%e:\ncan't parse dstart parameter - wrong input", err)
	}

	if repeat == "" {
		return "", fmt.Errorf("can't parse repeat parameter - empty input")
	}
	repeatData := strings.Split(repeat, " ")
	if len(repeatData) == 0 {
		return "", fmt.Errorf("can't parse repeat parameter - wrong input")
	}

	switch repeatData[0] {

	case "d":
		if len(repeatData) == 1 {
			return "", fmt.Errorf("can't shift date - no d's parameters on input")
		}
		shift, err := strconv.Atoi(repeatData[1])
		if err != nil {
			return "", fmt.Errorf("%e:\ncan't parse days shift parameter", err)
		}
		if shift > maxDays || shift < 0 {
			return "", fmt.Errorf("can't shift date - wrong input (%d) days", shift)
		}
		for {
			nextTime = nextTime.AddDate(0, 0, shift)
			if nextTime.After(now) {
				break
			}
		}

	case "y":
		for {
			nextTime = nextTime.AddDate(1, 0, 0)
			if nextTime.After(now) {
				break
			}
		}

	case "w":
		days := strings.Split(repeatData[1], ",")
		weekDays := make(map[int]bool, len(days))
		for _, d := range days {
			dayNum, err := strconv.Atoi(d)
			if err != nil {
				return "", fmt.Errorf("%e:\ncan't parse weekday parameter - unsupported format", err)
			}
			if dayNum < 1 || dayNum > 7 {
				return "", fmt.Errorf("can't use weekday parameter - wrong input")
			}
			if dayNum == 7 {
				dayNum = 0
			}
			weekDays[dayNum] = true
		}
		for {
			nextTime = nextTime.AddDate(0, 0, 1)
			if weekDays[int(nextTime.Weekday())] && nextTime.After(now) {
				break
			}
		}

	case "m":
		if len(repeatData) == 1 {
			return "", fmt.Errorf("can't shift date - no m's []<> parameters on input")
		}
		daysParse := strings.Split(repeatData[1], ",")
		days := make(map[int]bool, len(daysParse))
		daysTheseMonth := daysIn(int(nextTime.Month()), nextTime.Year())
		if len(repeatData) > 1 {
			for _, d := range daysParse {
				date, err := strconv.Atoi(d)
				if err != nil {
					return "", fmt.Errorf("%e:\ncan't parse m <date> parameter - unsupported format", err)
				}
				if date < -2 || date > 31 {
					return "", fmt.Errorf("can't use m <date> parameter - wrong input (%d)", date)
				}

				if date > daysTheseMonth {
					days[daysTheseMonth] = true
					continue
				}
				if date < 0 {
					days[daysTheseMonth+date+1] = true
					continue
				}
				days[date] = true
			}

			if len(repeatData) == 2 {
				for {
					nextTime = nextTime.AddDate(0, 0, 1)
					if days[int(nextTime.Day())] && nextTime.After(now) {
						break
					}
				}
				return nextTime.Format(layout), nil
			}

			if len(repeatData) == 3 {
				monthesParse := strings.Split(repeatData[2], ",")
				monthes := make(map[int]bool, len(monthesParse))

				for _, m := range monthesParse {
					month, err := strconv.Atoi(m)
					if err != nil {
						return "", fmt.Errorf("%e:\ncan't parse m [month] parameter - unsupported format", err)
					}
					if month < 0 || month > 12 {
						return "", fmt.Errorf("can't use m [month] parameter - wrong input (%d)", month)
					}
					monthes[month] = true
				}

				for {
					nextTime = nextTime.AddDate(0, 0, 1)
					if monthes[int(nextTime.Month())] && days[int(nextTime.Day())] && nextTime.After(now) {
						break
					}
				}
				return nextTime.Format(layout), nil
			}
		}

	default:
		return "", fmt.Errorf("can't use repeat parameter - undefined parameter (%q)", repeatData[0])
	}

	return nextTime.Format(layout), nil
}
