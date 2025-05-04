package main

import (
	"fmt"
)

type CompetitionSystem struct {
	Config      *Config
	EventLog    []string
	Competitors map[int]*Competitor
}

func NewCompetitionSystem(cfg *Config) *CompetitionSystem {
	return &CompetitionSystem{
		Config:      cfg,
		EventLog:    []string{},
		Competitors: make(map[int]*Competitor),
	}
}

func (cs *CompetitionSystem) HandleEvent(e Event) {
	c := cs.getOrCreateCompetitor(e.CompetitorID)
	logEntry := c.ProcessEvent(cs.Config, e)
	cs.EventLog = append(cs.EventLog, fmt.Sprintf("[%s] %s", e.Time.Format("15:04:05.000"), logEntry))
}

func (cs *CompetitionSystem) getOrCreateCompetitor(id int) *Competitor {
	if comp, ok := cs.Competitors[id]; ok {
		return comp
	}
	cs.Competitors[id] = NewCompetitor(id)
	return cs.Competitors[id]
}

func (cs *CompetitionSystem) PrintFinalReport() {
	for _, comp := range cs.Competitors {
		fmt.Println(comp.Report(cs.Config))
	}
}
