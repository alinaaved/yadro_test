package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	ExtraParams  []string
}

func ParseEvent(line string) (Event, error) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "[") {
		return Event{}, fmt.Errorf("invalid time format")
	}
	closeIdx := strings.Index(line, "]")
	if closeIdx == -1 {
		return Event{}, fmt.Errorf("missing closing bracket")
	}
	timeStr := line[1:closeIdx]
	eventTime, err := time.Parse("15:04:05.000", timeStr)
	if err != nil {
		return Event{}, fmt.Errorf("invalid time: %v", err)
	}

	rest := strings.TrimSpace(line[closeIdx+1:])
	parts := strings.Fields(rest)
	if len(parts) < 2 {
		return Event{}, fmt.Errorf("invalid event structure")
	}

	eventID, err := strconv.Atoi(parts[0])
	if err != nil {
		return Event{}, fmt.Errorf("invalid eventID")
	}

	competitorID, err := strconv.Atoi(parts[1])
	if err != nil {
		return Event{}, fmt.Errorf("invalid competitorID")
	}

	return Event{
		Time:         eventTime,
		EventID:      eventID,
		CompetitorID: competitorID,
		ExtraParams:  parts[2:],
	}, nil
}
