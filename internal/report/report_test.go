package report

import (
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"dial/internal/tracker"
)

// mirrors the real schema in internal/store/schema.sql, including the
// paused_at/paused_seconds columns the tracker package relies on.
const testSchema = `
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task TEXT NOT NULL,
    project TEXT,
    tags TEXT,
    started_at DATETIME NOT NULL,
    ended_at DATETIME,
    paused_at DATETIME,
    paused_seconds INTEGER DEFAULT 0,
    source TEXT DEFAULT 'manual',
    note TEXT
);
`

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if _, err := db.Exec(testSchema); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// Regression test: a session that is started and then stopped through the
// tracker must be visible to every report query used by the dashboard
// (TagBreakdown, DailyTotals, Summary) once it has ended. This previously
// broke because Summary -> currentStreak used SQLite's date() function on
// started_at, which is stored in Go's RFC3339 format (e.g.
// "2026-07-25T18:04:36.287292-04:00"). SQLite's date() can't parse the "T"
// separator and silently returns NULL, which then failed to scan into a
// non-nullable string and made the whole Summary() call error out. Since the
// frontend fetches TagBreakdown/DailyTotals/Summary with Promise.all, that
// one failure prevented the pie chart and bar chart from updating too, even
// though their own queries were fine.
func TestStoppedSession_VisibleInAllReports(t *testing.T) {
	db := newTestDB(t)
	trk := tracker.New(db)

	if _, err := trk.Start("write tests", "dial", "coding", "gui"); err != nil {
		t.Fatalf("Start error: %v", err)
	}
	time.Sleep(1100 * time.Millisecond)
	if _, err := trk.Stop(); err != nil {
		t.Fatalf("Stop error: %v", err)
	}

	start, end := RangeBounds("today", time.Now())

	tags, err := TagBreakdown(db, start, end)
	if err != nil {
		t.Fatalf("TagBreakdown error: %v", err)
	}
	if len(tags) != 1 || tags[0].Tag != "coding" || tags[0].Seconds < 1 {
		t.Errorf("TagBreakdown = %+v, want one entry for tag=coding with >=1 second", tags)
	}

	daily, err := DailyTotals(db, start, end)
	if err != nil {
		t.Fatalf("DailyTotals error: %v", err)
	}
	if len(daily) != 1 || daily[0].Seconds < 1 {
		t.Errorf("DailyTotals = %+v, want one day entry with >=1 second", daily)
	}

	stats, err := Summary(db, start, end)
	if err != nil {
		t.Fatalf("Summary error: %v (this is the bug that hid stopped sessions from the dashboard)", err)
	}
	if stats.SessionCount != 1 {
		t.Errorf("Summary.SessionCount = %d, want 1", stats.SessionCount)
	}
	if stats.TotalSeconds < 1 {
		t.Errorf("Summary.TotalSeconds = %d, want >=1", stats.TotalSeconds)
	}
	if stats.StreakDays < 1 {
		t.Errorf("Summary.StreakDays = %d, want >=1", stats.StreakDays)
	}
}

func TestCurrentStreak_HandlesISOFormattedTimestamps(t *testing.T) {
	db := newTestDB(t)
	trk := tracker.New(db)

	if _, err := trk.Start("task", "", "", "cli"); err != nil {
		t.Fatalf("Start error: %v", err)
	}
	if _, err := trk.Stop(); err != nil {
		t.Fatalf("Stop error: %v", err)
	}

	streak, err := currentStreak(db)
	if err != nil {
		t.Fatalf("currentStreak error: %v", err)
	}
	if streak != 1 {
		t.Errorf("currentStreak() = %d, want 1 for a session started today", streak)
	}
}

func TestTagBreakdown_UntaggedSessionsGroupedTogether(t *testing.T) {
	db := newTestDB(t)
	trk := tracker.New(db)

	if _, err := trk.Start("no tags task", "", "", "cli"); err != nil {
		t.Fatalf("Start error: %v", err)
	}
	if _, err := trk.Stop(); err != nil {
		t.Fatalf("Stop error: %v", err)
	}

	start, end := RangeBounds("today", time.Now())
	tags, err := TagBreakdown(db, start, end)
	if err != nil {
		t.Fatalf("TagBreakdown error: %v", err)
	}
	if len(tags) != 1 || tags[0].Tag != "untagged" {
		t.Errorf("TagBreakdown = %+v, want a single \"untagged\" entry", tags)
	}
}

func TestRangeBounds_TodayExcludesYesterday(t *testing.T) {
	db := newTestDB(t)
	trk := tracker.New(db)

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	if _, err := trk.Log("yesterday task", "", "work", "cli", yesterday, yesterday.Add(time.Hour)); err != nil {
		t.Fatalf("Log error: %v", err)
	}

	start, end := RangeBounds("today", now)
	tags, err := TagBreakdown(db, start, end)
	if err != nil {
		t.Fatalf("TagBreakdown error: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("TagBreakdown for today = %+v, want empty (session was logged yesterday)", tags)
	}
}
