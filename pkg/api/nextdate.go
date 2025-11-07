package api

import (
	"fmt"
	"gofinalproject/pkg/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var maxDays int = 400

// nextDayHandler handles Get-request with repeat pattern and responses next date.
func nextDayHandler(res http.ResponseWriter, req *http.Request) {
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	now := req.FormValue("now")
	if now == "" {
		now = time.Now().Format(db.Layout)
	}
	nowTime, err := time.Parse(db.Layout, now)
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
	_, _ = res.Write([]byte(next))
}

// NextDate calculates next date with different patterns, and dates.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	nextTime, err := time.Parse(db.Layout, dstart)
	if err != nil {
		return "", fmt.Errorf("can't parse start date parameter - wrong input; %e", err)
	}
	if repeat == "" {
		return "", fmt.Errorf("can't parse repeat parameter - empty input")
	}
	repeatData := strings.Fields(repeat)
	if len(repeatData) == 0 {
		return "", fmt.Errorf("can't parse repeat parameter - wrong input")
	}

	switch repeatData[0] {
	case "d":
		return dailyCalc(repeatData, now, nextTime)
	case "y":
		return yearlyCalc(now, nextTime)
	case "w":
		return weeklyCalc(repeatData, now, nextTime)
	case "m":
		return monthlyCalc(repeatData, now, nextTime)
	default:
		return "", fmt.Errorf("can't use repeat parameter - undefined parameter %q", repeatData[0])
	}
}

// dailyCalc calculates next date with daily pattern.
func dailyCalc(repeatData []string, now, nextTime time.Time) (string, error) {
	if len(repeatData) == 1 {
		return "", fmt.Errorf("can't shift date - no parameters in daily pattern")
	}
	shift, err := strconv.Atoi(repeatData[1])
	if err != nil {
		return "", fmt.Errorf("days parameter: unsupported format; %e", err)
	}
	if shift > maxDays || shift <= 0 {
		return "", fmt.Errorf("days parameter: wrong input (%d)", shift)
	}
	for {
		//wrong logic - now: 20240126 dstart: 20240202 repeat: d 30; nextTime: 20240303?
		nextTime = nextTime.AddDate(0, 0, shift)
		if isAfter(nextTime, now) {
			break
		}
	}
	return nextTime.Format(db.Layout), nil
}

// yearlyCalc calculates next date with yearly pattern.
func yearlyCalc(now, nextTime time.Time) (string, error) {
	for {
		nextTime = nextTime.AddDate(1, 0, 0)
		if isAfter(nextTime, now) {
			break
		}
	}
	return nextTime.Format(db.Layout), nil
}

// weeklyCalc calculates next date with weekly pattern. It builds a map with valid weekdays and then validate it in cycle.
func weeklyCalc(repeatData []string, now, nextTime time.Time) (string, error) {
	if len(repeatData) == 1 {
		return "", fmt.Errorf("can't shift date - no parameters in weekly pattern")
	}
	days := strings.Split(repeatData[1], ",")
	weekDays := make(map[int]bool, len(days))
	for _, d := range days {
		dayNum, err := strconv.Atoi(d)
		if err != nil {
			return "", fmt.Errorf("weekday parameter: unsupported format; %e", err)
		}
		if dayNum < 1 || dayNum > 7 {
			return "", fmt.Errorf("weekday parameter: wrong input (%d)", dayNum)
		}
		if dayNum == 7 {
			dayNum = 0
		}
		weekDays[dayNum] = true
	}
	for {
		dayNum := int(nextTime.Weekday())
		if weekDays[dayNum] && isAfter(nextTime, now) {
			break
		}
		nextTime = nextTime.AddDate(0, 0, 1)
	}
	return nextTime.Format(db.Layout), nil
}

// monthlyCalc calculates next date with monthly pattern. It uses a maps parsed monthlyParser
// and then validate it in cycles.
func monthlyCalc(repeatData []string, now, nextTime time.Time) (string, error) {
	daysCurrMonth := daysIn(int(nextTime.Month()), nextTime.Year())
	dayNums, monthNums, err := monthlyParser(repeatData, daysCurrMonth)
	if err != nil {
		return "", fmt.Errorf("monthlyParser error: %e", err)
	}

	if dayNums == nil && monthNums == nil {
		return "", fmt.Errorf("monthlyParser error: parse fail")
	}

	if dayNums != nil && monthNums == nil {
		for {
			dayNum := nextTime.Day()
			if dayNums[dayNum] && isAfter(nextTime, now) {
				break
			}
			nextTime = nextTime.AddDate(0, 0, 1)
		}
		return nextTime.Format(db.Layout), nil
	}

	for {
		dayNum := nextTime.Day()
		monthNum := nextTime.Month()
		if monthNums[int(monthNum)] && dayNums[dayNum] && isAfter(nextTime, now) {
			break
		}
		nextTime = nextTime.AddDate(0, 0, 1)
	}
	return nextTime.Format(db.Layout), nil

}

// monthlyParser parses and checks monthly string patterns and writes a maps with valid day numbers and month numbers
func monthlyParser(repeatData []string, daysThisMonth int) (map[int]bool, map[int]bool, error) {

	if len(repeatData) == 1 {
		return nil, nil, fmt.Errorf("can't shift date - no parameters in monthly pattern")
	}
	daysParse := strings.Split(repeatData[1], ",")
	dayNums := make(map[int]bool, len(daysParse))
	for _, d := range daysParse {
		date, err := strconv.Atoi(d)
		if err != nil {
			return nil, nil, fmt.Errorf("date parameter: unsupported format; %e", err)
		}
		if date < -2 || date == 0 || date > 31 {
			return nil, nil, fmt.Errorf("date parameter: wrong input (%d)", date)
		}
		if date < 0 {
			dayNums[daysThisMonth+date+1] = true
			continue
		}
		dayNums[date] = true
	}
	if len(repeatData) == 3 {
		monthesParse := strings.Split(repeatData[2], ",")
		monthNums := make(map[int]bool, len(monthesParse))

		for _, m := range monthesParse {
			month, err := strconv.Atoi(m)
			if err != nil {
				return nil, nil, fmt.Errorf("month parameter - unsupported format; %e", err)
			}
			if month < 0 || month > 12 {
				return nil, nil, fmt.Errorf("month parameter: wrong input (%d)", month)
			}
			monthNums[month] = true
		}
		return dayNums, monthNums, nil
	}
	return dayNums, nil, nil
}

// daysIn calculates the number of days in the passed month of passed year.
func daysIn(month, year int) int {
	t := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	return t.Day()
}

// isAfter returns true if date1 is after date2.
func isAfter(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	if y1 != y2 {
		return y1 > y2
	}
	if m1 != m2 {
		return m1 > m2
	}
	return d1 > d2
}
