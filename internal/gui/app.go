package gui

import (
	"context"
	"database/sql"
	"time"
	
	"dial/internal/report"
	"dial/internal/tracker"
)

type App struct {
	ctx context.Context
	db *sql.DB
	trk *tracker.Tracker
}

func NewApp(db *sql.DB) *App {
	return &App{db: db, trk: tracker.New(db)}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// SessionDTO is the JSON shape exposed to the frontend
type SessionDTO struct {
	ID int64 `json:"id"`
	Task string `json:"task"`
	Project string `json:"project"`
	Tags string `json:"tags"`
	StartedAt string `json:"startedAt"`
	EndedAt string `json:"endedAt,omitempty"`
	IsPaused bool `json:"isPaused"`
	ElapsedSecs int64 `json:"elapsedSeconds"`
}

type TagTimeDTO struct {
	Tag string `json:"tag"`
	Seconds int64 `json:"seconds"`
}

type DayTotalDTO struct {
	Date string `json:"date"`
	Seconds int64 `json:"seconds"`
}

type StatsDTO struct {
	TotalSeconds int64 `json:"totalSeconds"`
	SessionCount int `json:"sessionCount"`
	TopTag string `json:"topTag"`
	StreakDays int `json:"streakDays"`
}

func toDTO(s *tracker.Session) *SessionDTO {
	if s == nil {
		return nil
	}
	dto := &SessionDTO{
		ID: s.ID,
		Task: s.Task,
		Project: s.Project,
		Tags: s.Tags,
		StartedAt: s.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
		IsPaused: s.IsPaused(),
		ElapsedSecs: int64(s.Elapsed().Seconds()),
	}
	if s.EndedAt != nil {
		dto.EndedAt = s.EndedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	return dto
}

func (a *App) StartSession(task, project, tags string) (*SessionDTO, error) {
	s, err := a.trk.Start(task, project, tags, "gui")
	if err != nil {
		return nil, err
	}
	return toDTO(s), nil
}

func (a *App) StopSession() (*SessionDTO, error) {
	s, err := a.trk.Stop()
	if err != nil {
		return nil, err
	}
	return toDTO(s), nil
}

func (a *App) PauseSession() (*SessionDTO, error) {
	s, err := a.trk.Pause()
	if err != nil {
		return nil, err
	}
	return toDTO(s), nil
}

func (a *App) ResumeSession() (*SessionDTO, error) {
	s, err := a.trk.Resume()
	if err != nil {
		return nil, err
	}
	return toDTO(s), nil
}

func (a *App) CurrentSession() (*SessionDTO, error) {
	s, err := a.trk.Current()
	if err != nil {
		if err == tracker.ErrNoRunningSession {
			return nil, nil
		}
		return nil, err
	}
	return toDTO(s), nil
}

func (a *App) GetTagBreakdown(rangeName string) ([]TagTimeDTO, error) {
	start, end := report.RangeBounds(rangeName, time.Now())
	rows, err := report.TagBreakdown(a.db, start, end)
	if err != nil {
		return nil, err
	}
	out := make([]TagTimeDTO, len(rows))
	for i, r := range rows {
		out[i] = TagTimeDTO{Tag: r.Tag, Seconds: r.Seconds}
	}
	return out, nil
}

// Returns time-per-day for a named range
func (a *App) GetDailyTotals(rangeName string) ([]DayTotalDTO, error) {
	start, end := report.RangeBounds(rangeName, time.Now())
	rows, err := report.DailyTotals(a.db, start, end)
	if err != nil {
		return nil, err
	}
	out := make([]DayTotalDTO, len(rows))
	for i, r := range rows {
		out[i] = DayTotalDTO{Date: r.Date, Seconds: r.Seconds}
	}
	return out, nil
}

// Returns summary stats for a named range
func (a *App) GetStats(rangeName string) (*StatsDTO, error) {
	start, end := report.RangeBounds(rangeName, time.Now())
	s, err := report.Summary(a.db, start, end)
	if err != nil {
		return nil, err
	}
	return &StatsDTO{
		TotalSeconds: s.TotalSeconds,
		SessionCount: s.SessionCount,
		TopTag: s.TopTag,
		StreakDays: s.StreakDays
	}, nil
}
