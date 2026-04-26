package db

import (
	"fmt"
)

// HourlyBucket holds command count for one hour of the day.
type HourlyBucket struct {
	Hour  int
	Count int64
}

// DailyBucket holds command count for one calendar day.
type DailyBucket struct {
	Date  string
	Count int64
}

// ProjectSummary is a full per-project stats row used by `pulse projects`.
type ProjectSummary struct {
	Name        string
	TotalTimeMS int64
	Commands    int64
	SuccessRate float64
}

func (db *DB) GetHourlyStats(date string) ([]HourlyBucket, error) {
	start, end, err := localDateRange(date)
	if err != nil {
		return nil, fmt.Errorf("parsing date %q: %w", date, err)
	}
	rows, err := db.conn.Query(`
		SELECT CAST(strftime('%H', created_at, 'localtime') AS INTEGER) AS hr, COUNT(*) AS cnt
		FROM commands
		WHERE created_at >= ? AND created_at < ?
		GROUP BY hr
		ORDER BY hr`, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []HourlyBucket
	for rows.Next() {
		var b HourlyBucket
		if err := rows.Scan(&b.Hour, &b.Count); err != nil {
			continue
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// GetDailyActivity returns command count per day for the last N days, oldest first.
func (db *DB) GetDailyActivity(days int) ([]DailyBucket, error) {
	since := localWindowStart(days)
	rows, err := db.conn.Query(`
		SELECT DATE(created_at, 'localtime') AS day, COUNT(*) AS cnt
		FROM commands
		WHERE created_at >= ?
		GROUP BY day
		ORDER BY day ASC`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []DailyBucket
	for rows.Next() {
		var b DailyBucket
		if err := rows.Scan(&b.Date, &b.Count); err != nil {
			continue
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// GetTodayStats returns a Stats summary scoped to today only.
func (db *DB) GetTodayStats() (*Stats, error) {
	start, end := localDayRange(nowFunc())
	var s Stats
	err := db.conn.QueryRow(`
		SELECT
			COALESCE(SUM(CASE WHEN noise = 0 THEN 1 ELSE 0 END), 0)  AS dev_total,
			COALESCE(SUM(CASE WHEN noise = 1 THEN 1 ELSE 0 END), 0)  AS noise_total,
			COALESCE(SUM(CASE WHEN noise = 0 THEN duration_ms ELSE 0 END), 0) AS total_ms,
			COALESCE(AVG(CASE WHEN noise = 0 AND exit_code = 0 THEN 1.0
			               WHEN noise = 0 THEN 0.0 END) * 100, 0)   AS success_rate
		FROM commands
		WHERE created_at >= ? AND created_at < ?`, start, end).
		Scan(&s.TotalCommands, &s.NoiseCommands, &s.TotalTimeMS, &s.SuccessRate)
	return &s, err
}

// GetProjectList returns all projects with full stats for the last N days.
func (db *DB) GetProjectList(days int) ([]ProjectSummary, error) {
	since := localWindowStart(days)
	rows, err := db.conn.Query(`
		SELECT
			project,
			COALESCE(SUM(CASE WHEN noise = 0 THEN duration_ms ELSE 0 END), 0) AS total_ms,
			COALESCE(SUM(CASE WHEN noise = 0 THEN 1 ELSE 0 END), 0) AS total_cmds,
			COALESCE(AVG(CASE WHEN noise = 0 AND exit_code = 0 THEN 1.0
			                  WHEN noise = 0 THEN 0.0 END) * 100, 0) AS success_rate
		FROM commands
		WHERE created_at >= ? AND project != ''
		GROUP BY project
		HAVING total_cmds > 0
		ORDER BY total_ms DESC`, since)
	if err != nil {
		return nil, fmt.Errorf("querying project list: %w", err)
	}
	defer rows.Close()

	var out []ProjectSummary
	for rows.Next() {
		var p ProjectSummary
		if err := rows.Scan(&p.Name, &p.TotalTimeMS, &p.Commands, &p.SuccessRate); err != nil {
			continue
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
