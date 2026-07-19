package tracker

import (
    "database/sql"
    "errors"
    "time"
)

type Session struct {
    ID int64
    Task string
    Project string
    Tags string
    StartedAt time.Time
    EndedAt *time.Time
	PausedAt *time.Time
	PausedSeconds int64
    Source string
    Note string
}

// Elapsed returns time actually worked, exclusing paused time
// If the session is currently paused, time since PausedAt doesnt count
func (s *Session) Elapsed() time.Duration {
	end := time.Now()
	if s.EndedAt != nil {
		end = *s.EndedAt
	} else if s.PausedAt != nil {
		end = *s.PausedAt
	}

	total := end.Sub(s.StartedAt) - time.Duration(s.PausedSeconds)*time.Second
	if total < 0 {
		total = 0
	}
	return total
}

func (s *Session) IsPaused() bool {
	return s.PausedAt != nil && s.EndedAt == nil
}

type Tracker struct {
    db *sql.DB
}

func New(db *sql.DB) *Tracker {
    return &Tracker{db: db}
}

var ErrNoRunningSession = errors.New("no running session")

// start stops any currently running session, then starts a new one
func (t *Tracker) Start(task, project, tags, source string) (*Session, error) {
    if _, err := t.Stop(); err != nil && !errors.Is(err, ErrNoRunningSession) {
        return nil, err
    }

    now := time.Now()
    res, err := t.db.Exec(
        `INSERT INTO sessions (task, project, tags, started_at, source) VALUES (?, ?, ?, ?, ?)`,
	task, project, tags, now, source,
    )
    if err != nil {
        return nil, err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }

    return &Session{
        ID: id,
	Task: task,
	Project: project,
	Tags: tags,
	StartedAt: now,
	Source: source,
    }, nil
}

// stop ends the current running session, if any.
func (t *Tracker) Stop() (*Session, error) {
    row := t.db.QueryRow(
        `SELECT id, task, project, tags, started_at, source
	 	 FROM sessions WHERE ended_at IS NULL
	 	 ORDER BY started_at DESC LIMIT 1`,
    )

    var s Session
    if err := row.Scan(&s.ID, &s.Task, &s.Project, &s.Tags, &s.StartedAt, &s.PausedAt, &s.PausedSeconds, &s.Source); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNoRunningSession
	}
	return nil, err
    }

    now := time.Now()
	// if still paused when stopped, fold the final
	// pause span into paused_seconds
	if s.PausedAt != nil {
		s.PausedSeconds += int64(now.Sub(*s.PausedAt).Seconds())
	}
    if _, err := t.db.Exec(`UPDATE sessions SET ended_at = ?, paused_at = NULL, paused_seconds = ? WHERE id = ?`, now, s.PausedSeconds, s.ID); err != nil {
        return nil, err
    }
    s.EndedAt = &now
	s.PausedAt = nil 
    return &s, nil
}

func (t *Tracker) Pause() (*Session, error) {
	s, err := t.Current()
	if err != nil {
		return nil, err
	}
	if s.PausedAt != nil {
		return s, nil // already paused, no-op
	}
	now := time.Now()
	if _, err := t.db.Exec(`UPDATE sessions SET paused_at = ? WHERE id = ?`, now, s.ID); err != nil {
		return nil, err
	}
	s.PausedAt = &now
	return s, nil
}

func (t *Tracker) Resume() (*Session, error) {
	s, err := t.Current()
	if err != nil {
		return nil, err
	}
	if s.PausedAt == nil {
		return s, nil // not paused, no-op
	}
	now := time.Now()
	pausedSpan := int64(now.Sub(*s.PausedAt).Seconds())
	newPausedSeconds := s.PausedSeconds + pausedSpan
	if _, err := t.db.Exec(
		`UPDATE sessions SET paused_at = NULL, paused_seconds = ? WHERE id = ?`,
		newPausedSeconds, s.ID,
	); err != nil {
		return nil, err
	}
	s.PausedAt = nil
	s.PausedSeconds = newPausedSeconds
	return s, nil
}

// returns the currently running session
func (t *Tracker) Current() (*Session, error) {
    row := t.db.QueryRow(
        `SELECT id, task, project, tags, started_at, paused_at, paused_seconds, source
	 FROM sessions WHERE ended_at IS NULL
	 ORDER BY started_at DESC LIMIT 1`,
    )

    var s Session
    if err := row.Scan(&s.ID, &s.Task, &s.Project, &s.Tags, &s.StartedAt, &s.PausedAt, &s.PausedSeconds, &s.Source); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
	    	return nil, ErrNoRunningSession
		}
		return nil, err
    }
    return &s, nil
}

// Log inserts a completed session
func (t *Tracker) Log(task, project, tags, source string, start, end time.Time) (*Session, error) {
    res, err := t.db.Exec(
	`INSERT INTO sessions (task, project, tags, started_at, ended_at, source) VALUES (?, ?, ?, ?, ?, ?)`,
	task, project, tags, start, end, source,
    )
    if err != nil {
        return nil, err
    }
    id, err := res.LastInsertId()
    if err != nil {
	return nil, err
    }
    return &Session{
	ID: id, Task: task, Project: project, Tags: tags,
	StartedAt: start, EndedAt: &end, Source: source,
    }, nil
}
