package tracker

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

const testSchema = `
CREATE TABLE sessions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	task TEXT NOT NULL,
	project TEXT,
	tags TEXT,
	started_at DATETIME NOT NULL,
	ended_at DATETIME,
	source TEXT DEFAULT 'manual',
	note TEXT
);
`

func newTestTracker(t *testing.T) *Tracker {
	t.Helper()
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if _, err := db.Exec(testSchema); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return New(db)
}

func TestStart_CreatesRunningSession(t *testing.T) {
	trk := newTestTracker(t)

	s, err := trk.Start("write tests", "dial", "coding", "cli")
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	if s.Task != "write tests" {
		t.Errorf("Task = %q, want %q", s.Task, "write tests")
	}
	if s.EndedAt != nil {
		t.Errorf("EndedAt = %v, want nil for a freshly started session", s.EndedAt)
	}
	if s.ID == 0 {
		t.Error("expected non-zero ID after insert")
	}
}

func TestStart_StopsPreviousRunningSession(t *testing.T) {
	trk := newTestTracker(t)

	first, err := trk.Start("task one", "", "", "cli")
	if err != nil {
		t.Fatalf("Start(first) error: %v", err)
	}

	if _, err := trk.Start("task two", "", "", "cli"); err != nil {
		t.Fatalf("Start(second) error: %v", err)
	}

	// the first session should now be closed out
	row := trk.db.QueryRow(`SELECT ended_at FROM sessions WHERE id = ?`, first.ID)
	var ended *time.Time
	if err := row.Scan(&ended); err != nil {
		t.Fatalf("scan ended_at: %v", err)
	}
	if ended == nil {
		t.Error("expected first session to have ended_at set after starting a second session")
	}

	// only one session should be currently running
	cur, err := trk.Current()
	if err != nil {
		t.Fatalf("Current() error: %v", err)
	}
	if cur.Task != "task two" {
		t.Errorf("Current().Task = %q, want %q", cur.Task, "task two")
	}
}

func TestStop_NoRunningSession(t *testing.T) {
	trk := newTestTracker(t)

	_, err := trk.Stop()
	if !errors.Is(err, ErrNoRunningSession) {
		t.Errorf("Stop() error = %v, want ErrNoRunningSession", err)
	}
}

func TestStop_EndsRunningSession(t *testing.T) {
	trk := newTestTracker(t)

	started, err := trk.Start("focus block", "", "", "cli")
	if err != nil {
		t.Fatalf("Start error: %v", err)
	}

	stopped, err := trk.Stop()
	if err != nil {
		t.Fatalf("Stop error: %v", err)
	}
	if stopped.ID != started.ID {
		t.Errorf("Stop() returned session ID %d, want %d", stopped.ID, started.ID)
	}
	if stopped.EndedAt == nil {
		t.Fatal("expected EndedAt to be set after Stop")
	}
	if stopped.EndedAt.Before(stopped.StartedAt) {
		t.Error("EndedAt should not be before StartedAt")
	}

	// calling stop again should now report no running session
	if _, err := trk.Stop(); !errors.Is(err, ErrNoRunningSession) {
		t.Errorf("second Stop() error = %v, want ErrNoRunningSession", err)
	}
}

func TestCurrent_NoRunningSession(t *testing.T) {
	trk := newTestTracker(t)

	_, err := trk.Current()
	if !errors.Is(err, ErrNoRunningSession) {
		t.Errorf("Current() error = %v, want ErrNoRunningSession", err)
	}
}

func TestCurrent_ReturnsMostRecentRunning(t *testing.T) {
	trk := newTestTracker(t)

	if _, err := trk.Start("task", "proj-a", "tag-a", "cli"); err != nil {
		t.Fatalf("Start error: %v", err)
	}

	cur, err := trk.Current()
	if err != nil {
		t.Fatalf("Current error: %v", err)
	}
	if cur.Project != "proj-a" || cur.Tags != "tag-a" {
		t.Errorf("Current() = %+v, want project=proj-a tags=tag-a", cur)
	}
}

func TestLog_InsertsCompletedSession(t *testing.T) {
	trk := newTestTracker(t)

	end := time.Now()
	start := end.Add(-45 * time.Minute)

	s, err := trk.Log("retroactive task", "dial", "planning", "cli", start, end)
	if err != nil {
		t.Fatalf("Log error: %v", err)
	}
	if s.EndedAt == nil {
		t.Fatal("expected EndedAt to be set on a logged session")
	}
	if !s.EndedAt.Equal(end) {
		t.Errorf("EndedAt = %v, want %v", s.EndedAt, end)
	}
	if s.StartedAt.After(*s.EndedAt) {
		t.Error("StartedAt should not be after EndedAt")
	}

	// log should not affect currently running session state
	if _, err := trk.Current(); !errors.Is(err, ErrNoRunningSession) {
		t.Errorf("Current() error = %v, want ErrNoRunningSession after Log", err)
	}
}

func TestLog_DoesNotInterfereWithRunningSession(t *testing.T) {
	trk := newTestTracker(t)

	running, err := trk.Start("ongoing", "", "", "cli")
	if err != nil {
		t.Fatalf("Start error: %v", err)
	}

	end := time.Now().Add(-time.Hour)
	start := end.Add(-30 * time.Minute)
	if _, err := trk.Log("past task", "", "", "cli", start, end); err != nil {
		t.Fatalf("Log error: %v", err)
	}

	cur, err := trk.Current()
	if err != nil {
		t.Fatalf("Current error: %v", err)
	}
	if cur.ID != running.ID {
		t.Errorf("Current().ID = %d, want %d (running session should be untouched)", cur.ID, running.ID)
	}
}
