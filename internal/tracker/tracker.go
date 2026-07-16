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
    Ended At *time.Time
    Source string
    Note string
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
    if err := row.Scan(&s.ID, &s.Task, &s.Project, &s.Tags, &s.StartedAt, &s.Source); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNoRunningSession
	}
	return nil, err
    }

    now := time.Now()
    if _, err := t.db.Exec(`UPDATE sessions SET ended_at = ? WHERE id = ?`, now, s.ID); err != nil {
        return nil, err
    }
    s.EndedAt = &now
    return &s, nil
}

// returns the currently running session
func (t *Tracker) Current() (*Session, error) {
    row := t.db.QueryRow(
        `SELECT id, task, project, tags, started_at, source
	 FROM sessions WHERE ended_at IS NULL
	 ORDER BY started_at DESC LIMIT 1`,
    )

    var s Session
    if err := row.Scan(&s.ID, &s.Task, &s.Project, &s.Tags, &s.StartedAt, &s.Source); err != nil {
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
