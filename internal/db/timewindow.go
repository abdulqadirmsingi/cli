package db

import "time"

var nowFunc = time.Now

const sqliteTimeLayout = "2006-01-02 15:04:05"

func sqliteUTC(t time.Time) string {
	return t.UTC().Format(sqliteTimeLayout)
}

func localDayRange(day time.Time) (string, string) {
	y, m, d := day.In(time.Local).Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 0, 1)
	return sqliteUTC(start), sqliteUTC(end)
}

func localDateRange(date string) (string, string, error) {
	day, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return "", "", err
	}
	start, end := localDayRange(day)
	return start, end, nil
}

func localWindowStart(days int) string {
	if days <= 1 {
		start, _ := localDayRange(nowFunc())
		return start
	}
	y, m, d := nowFunc().In(time.Local).Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	return sqliteUTC(today.AddDate(0, 0, -(days - 1)))
}
