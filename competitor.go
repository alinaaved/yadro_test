package main

import (
	"fmt"
	"strings"
	"time"
)

type Competitor struct {
	ID           int
	Registered   bool
	StartPlanned time.Time
	StartedAt    *time.Time
	Finished     bool
	NotStarted   bool
	NotFinished  bool
	Comment      string

	LapTimes    []time.Duration
	PenaltyTime time.Duration
	Shots       int
	Hits        int

	currentLapStart time.Time
	penaltyStart    *time.Time
}

func NewCompetitor(id int) *Competitor {
	return &Competitor{ID: id}
}

func (c *Competitor) ProcessEvent(cfg *Config, e Event) string {
	switch e.EventID {
	case 1:
		c.Registered = true
		return fmt.Sprintf("The competitor(%d) registered", c.ID)
	case 2:
		parsed, _ := time.Parse("15:04:05.000", e.ExtraParams[0])
		c.StartPlanned = parsed
		return fmt.Sprintf("The start time for the competitor(%d) was set by a draw to %s", c.ID, parsed.Format("15:04:05.000"))
	case 3:
		return fmt.Sprintf("The competitor(%d) is on the start line", c.ID)
	case 4:
		c.StartedAt = &e.Time
		c.currentLapStart = e.Time
		return fmt.Sprintf("The competitor(%d) has started", c.ID)
	case 5:
		line := e.ExtraParams[0]
		return fmt.Sprintf("The competitor(%d) is on the firing range(%s)", c.ID, line)
	case 6:
		c.Shots++
		c.Hits++
		return fmt.Sprintf("The target(%s) has been hit by competitor(%d)", e.ExtraParams[0], c.ID)
	case 7:
		return fmt.Sprintf("The competitor(%d) left the firing range", c.ID)
	case 8:
		c.penaltyStart = &e.Time
		return fmt.Sprintf("The competitor(%d) entered the penalty laps", c.ID)
	case 9:
		if c.penaltyStart != nil {
			c.PenaltyTime += e.Time.Sub(*c.penaltyStart)
			c.penaltyStart = nil
		}
		return fmt.Sprintf("The competitor(%d) left the penalty laps", c.ID)
	case 10:
		if !c.currentLapStart.IsZero() {
			c.LapTimes = append(c.LapTimes, e.Time.Sub(c.currentLapStart))
			c.currentLapStart = e.Time
		}
		return fmt.Sprintf("The competitor(%d) ended the main lap", c.ID)
	case 11:
		c.NotFinished = true
		if len(e.ExtraParams) > 0 {
			c.Comment = strings.Join(e.ExtraParams, " ")
		}
		return fmt.Sprintf("The competitor(%d) can`t continue: %s", c.ID, c.Comment)
	case 33:
		c.Finished = true
		return fmt.Sprintf("The competitor(%d) has finished", c.ID)
	case 32:
		c.NotStarted = true
		return fmt.Sprintf("The competitor(%d) is disqualified", c.ID)
	default:
		return fmt.Sprintf("Unknown event %d for competitor(%d)", e.EventID, c.ID)
	}
}

func (c *Competitor) Report(cfg *Config) string {
	status := ""
	if c.NotStarted {
		status = "[NotStarted]"
	} else if c.NotFinished {
		status = "[NotFinished]"
	} else if c.StartedAt != nil && len(c.LapTimes) == cfg.Laps {
		total := time.Duration(0)
		laps := []string{}
		for _, lap := range c.LapTimes {
			total += lap
			speed := float64(cfg.LapLen) / lap.Seconds()
			laps = append(laps, fmt.Sprintf("{%s, %.3f}", formatDur(lap), speed))
		}
		speedPenalty := float64(cfg.PenaltyLen) / c.PenaltyTime.Seconds()
		return fmt.Sprintf("%s %d [%s] {%s, %.3f} %d/%d", status, c.ID, strings.Join(laps, ", "), formatDur(c.PenaltyTime), speedPenalty, c.Hits, c.Shots)
	}
	return fmt.Sprintf("%s %d [] {,} {,} %d/%d", status, c.ID, c.Hits, c.Shots)
}

func formatDur(d time.Duration) string {
	return d.Truncate(time.Millisecond).String()
}
