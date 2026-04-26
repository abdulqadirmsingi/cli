package db

import (
	"path/filepath"
	"testing"
	"time"
)

func TestTodayStatsAndHourlyUseLocalCalendarDay(t *testing.T) {
	database := openTestDB(t)
	defer database.Close()

	fixedNow := time.Date(2026, 4, 26, 8, 30, 0, 0, time.Local)
	withNow(t, fixedNow)

	insertCommandAt(t, database, "go test ./...", "cli", 0, 90_000, false, localAt(2026, 4, 26, 8, 10))
	insertCommandAt(t, database, "ls", "cli", 0, 10_000, true, localAt(2026, 4, 26, 8, 12))
	insertCommandAt(t, database, "npm test", "web", 1, 30_000, false, localAt(2026, 4, 25, 23, 59))
	insertCommandAt(t, database, "go build", "cli", 0, 45_000, false, localAt(2026, 4, 27, 0, 1))

	stats, err := database.GetTodayStats()
	if err != nil {
		t.Fatalf("GetTodayStats: %v", err)
	}
	if stats.TotalCommands != 1 || stats.NoiseCommands != 1 {
		t.Fatalf("today command counts = dev %d noise %d, want dev 1 noise 1", stats.TotalCommands, stats.NoiseCommands)
	}
	if stats.TotalTimeMS != 90_000 {
		t.Fatalf("today total time = %d, want 90000", stats.TotalTimeMS)
	}

	hourly, err := database.GetHourlyStats("2026-04-26")
	if err != nil {
		t.Fatalf("GetHourlyStats: %v", err)
	}
	if got := bucketCount(hourly, 8); got != 2 {
		t.Fatalf("hour 8 count = %d, want 2", got)
	}
	if got := bucketCount(hourly, 23); got != 0 {
		t.Fatalf("hour 23 count = %d, want 0", got)
	}
}

func TestTodayTopSectionsDoNotUseRollingWindow(t *testing.T) {
	database := openTestDB(t)
	defer database.Close()

	fixedNow := time.Date(2026, 4, 26, 8, 30, 0, 0, time.Local)
	withNow(t, fixedNow)

	insertCommandAt(t, database, "go test", "cli", 0, 60_000, false, localAt(2026, 4, 26, 8, 10))
	insertCommandAt(t, database, "go build", "cli", 0, 30_000, false, localAt(2026, 4, 26, 8, 15))
	insertCommandAt(t, database, "docker build", "infra", 0, 900_000, false, localAt(2026, 4, 25, 22, 0))

	cmds, err := database.GetTodayTopCommands(3)
	if err != nil {
		t.Fatalf("GetTodayTopCommands: %v", err)
	}
	if len(cmds) != 1 || cmds[0].Name != "go" || cmds[0].Count != 2 {
		t.Fatalf("today top commands = %#v, want go x2 only", cmds)
	}

	projects, err := database.GetTodayTopProjects(3)
	if err != nil {
		t.Fatalf("GetTodayTopProjects: %v", err)
	}
	if len(projects) != 1 || projects[0].Name != "cli" || projects[0].MS != 90_000 {
		t.Fatalf("today top projects = %#v, want cli 90000ms only", projects)
	}
}

func TestProjectAggregatesMatchDevTimeSemantics(t *testing.T) {
	database := openTestDB(t)
	defer database.Close()

	fixedNow := time.Date(2026, 4, 26, 8, 30, 0, 0, time.Local)
	withNow(t, fixedNow)

	insertCommandAt(t, database, "go test", "cli", 0, 120_000, false, localAt(2026, 4, 26, 8, 10))
	insertCommandAt(t, database, "go test ./bad", "cli", 1, 60_000, false, localAt(2026, 4, 26, 8, 15))
	insertCommandAt(t, database, "ls", "cli", 0, 600_000, true, localAt(2026, 4, 26, 8, 20))
	insertCommandAt(t, database, "pwd", "noise-only", 0, 900_000, true, localAt(2026, 4, 26, 8, 25))

	projects, err := database.GetProjectList(1)
	if err != nil {
		t.Fatalf("GetProjectList: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("projects = %#v, want only cli", projects)
	}
	if projects[0].Name != "cli" || projects[0].TotalTimeMS != 180_000 || projects[0].Commands != 2 {
		t.Fatalf("cli project = %#v, want 180000ms and 2 commands", projects[0])
	}
	if projects[0].SuccessRate != 50 {
		t.Fatalf("cli success rate = %.1f, want 50.0", projects[0].SuccessRate)
	}

	top, err := database.GetTopProjects(1, 3)
	if err != nil {
		t.Fatalf("GetTopProjects: %v", err)
	}
	if len(top) != 1 || top[0].Name != "cli" || top[0].MS != 180_000 || top[0].Count != 2 {
		t.Fatalf("top projects = %#v, want cli with dev-only totals", top)
	}
}

func openTestDB(t *testing.T) *DB {
	t.Helper()
	database, err := Open(filepath.Join(t.TempDir(), "pulse.db"))
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	return database
}

func withNow(t *testing.T, now time.Time) {
	t.Helper()
	original := nowFunc
	nowFunc = func() time.Time { return now }
	t.Cleanup(func() { nowFunc = original })
}

func localAt(year int, month time.Month, day, hour, min int) time.Time {
	return time.Date(year, month, day, hour, min, 0, 0, time.Local)
}

func insertCommandAt(t *testing.T, database *DB, command, project string, exitCode int, durationMS int64, noise bool, at time.Time) {
	t.Helper()
	noiseInt := 0
	if noise {
		noiseInt = 1
	}
	_, err := database.conn.Exec(`
		INSERT INTO commands (command, directory, project, exit_code, duration_ms, noise, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		command, "/tmp/"+project, project, exitCode, durationMS, noiseInt, sqliteUTC(at))
	if err != nil {
		t.Fatalf("insert command %q: %v", command, err)
	}
}

func bucketCount(buckets []HourlyBucket, hour int) int64 {
	for _, b := range buckets {
		if b.Hour == hour {
			return b.Count
		}
	}
	return 0
}
