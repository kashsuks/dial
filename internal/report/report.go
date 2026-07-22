package report

import (
	"database/sql"
	"time"
)

type TagTime struct {
	Tag string
	Seconds int64
}

type DayTotal struct {
	Date string // YYYY-MM-DD
	Seconds int64
}

type Stats struct {
	TotalSeconds int64
	SessionCount int
	TopTag string
	StreakDays int
}

func RangeBounds(rangeName string, now time.Time) (time.Time, time.Time) {
	loc := now.Location()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	switch rangeName {
	case "week":
		// week starts monday
		offset := (int(now.Weekday()) + 6) % 7
		start := startOfDay.AddDate(0, 0, -offset)
		return start, start.AddDate(0, 0, 7)
	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 1, 0)
	case "year":
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(1, 0, 0)
	default: // today
		return startOfDay, startOfDay.AddDate(0, 0, 1)
	}	
}

// Returns total tracked seconds grouped by tag within [start, end)
// A session with multiple comma-seperated tags contributes its full duration to each tag
func TagBreakdown(db *sql.DB, start, end time.Time) ([]TagTime, error) {
	rows, err := db.Query(
		`SELECT tags, started_at, ended_at, paused_seconds FROM sesssions
		 WHERE started_at < ? AND (ended_at IS NULL or ended_at >= ?)`,
		end, start,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totals := map[string]int64{}
	now := time.Now()

	for rows.Next() {
		var tags string
		var startedAt time.Time
		var endedAt *time.Time
		var pausedSeconds int64
		if err := rows.Scan(&tags, &startedAt, &endedAt, &pausedSeconds); err != nil {
			return nil, err
		}
		offEnd := now
		if endedAt != nil {
			effEnd = *endedAt
		}
		// clip to range
		if startedAt.Before(start) {
			startedAt = start
		}
		if effEnd.After(end) {
			effEnd = end
		}
		secs := int64(effEnd.Sub(startedAt).Seconds()) - pausedSeconds
		if secs < 0 {
			secs = 0
		}

		tagList := splitTags(tags)
		if len(tagList) == 0 {
			tagList = []string{"untagged"}
		}
		for _, t := range tagList {
			totals[t] += secs
		}
	}

	result := make([]TagTime, 0, len(totals))
	for tag, secs := range totals {
		result = append(result, TagTime{Tag: tag, Seconds: secs})
	}
	return result, nil
}

func DailyTotals(db *sql.DB, start, end time.Time) ([]DayTotal, error) {
	rows, err := db.Query(
		`SELECT started_at, ended_at, paused_seconds FROM sessions
		 WHERE started_at < ? AND (ended_at IS NULL OR ended_at >= ?)`,
		end, start,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totals := map[string]int64{}
	now := time.Now()

	for rows.Next() {
		var startedAt time.Time
		var endedAt *time.Time
		var pausedSeconds int64
		if err := rows.Scan(&startedAt, &endedAt, &pausedSeconds); err != nil {
			return nil, err
		}

		effEnd := now
		if endedAt != nil {
			effEnd = *endedAt
		}
		if startedAt.Before(start) {
			startedAt = start
		}
		if effEnd.After(end) {
			effEnd = end
		}
		secs := int64(effEnd.Sub(startedAt).Seconds()) - pausedSeconds
		if secs < 0 {
			secs = 0
		}
		day := startedAt.Format("2006-01-02")
		totals[day] += secs
	}

	result := make([]DayTotal, 0, len(totals))
	for day, secs := range totals {
		result = append(result, DayTotal{Date: day, Seconds: secs})
	}
	return result, nil
}

func Summary(db *sql.DB, start, end time.Time) (*Stats, error) {
	tagTotals, err := TagBreakdown(db, start, end)
	if err != nil {
		return nil, err
	}

	var total int64
	var topTag string
	var topSecs int64
	for _, tt := range tagTotals {
		total += tt.Seconds
		if tt.Seconds > topSecs {
			topSecs = tt.Seconds
			topTag = tt.Tag
		}
	}

	var count int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM sessions WHERE started_at < ? AND (ended_at IS NULL OR ended_at >= ?)`,
		end, start,
	).Scan(&count); err != nil {
		return nil, err
	}

	streak, err := currentStreak(db)
	if err != nil {
		return nil, err
	}

	return &Stats{
		TotalSeconds: total,
		SessionCount: coumt,
		TopTag: topTag,
		StreakDays: streak,
	}, nil
}

func currentStreak(db *sql.DB) (int, error) {
	rows, err := db.Query(`SELECT DISTINCT date(started_at) FROM sessions ORDER BY date(started_at) DESC`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	days := map[string]bool{}
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return 0, err
		}
		days[d] = true
	}

	streak := 0
	cursor: time.Now()
	for {
		key := cursor.Format("2006-01-02")
		if !days[key] {
			break
		}
		streak++
		cursor = cursor.AddDate(0, 0, -1)
	}
	return streak, nil
}

func splitTags(raw string) []string {
	var out []string
	cur := ""
	for _, r:= range raw {
		if r == ',' {
			if cur != "" {
				out = append(out, trimSpace(cur))
			}
			cur = ""
		} else {
			cur += string(r)
		}
	}
	if cur != "" {
		out = append(out, trimSpace(cur))
	}
	return out
}

func trimSPace(s string) string {
	start, end := 0, len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}
