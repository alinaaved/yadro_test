package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Laps          int           `json:"laps"`
	LapLen        int           `json:"lapLen"`
	PenaltyLen    int           `json:"penaltyLen"`
	FiringLines   int           `json:"firingLines"`
	Start         time.Time     `json:"-"`
	StartRaw      string        `json:"start"`
	StartDelta    time.Duration `json:"-"`
	StartDeltaRaw string        `json:"startDelta"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	//обработка времени
	cfg.Start, err = time.Parse("15:04:05.000", cfg.StartRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %v", err)
	}

	cfg.StartDelta, err = time.ParseDuration(cfg.StartDeltaRaw + "s")
	if err != nil {
		cfg.StartDelta, err = parseHHMMSS(cfg.StartDeltaRaw)
		if err != nil {
			return nil, fmt.Errorf("invalid start delta: %v", err)
		}
	}

	return &cfg, nil
}

func parseHHMMSS(s string) (time.Duration, error) {
	parsed, err := time.Parse("15:04:05", s)
	if err != nil {
		return 0, err
	}
	return time.Duration(parsed.Hour())*time.Hour +
		time.Duration(parsed.Minute())*time.Minute +
		time.Duration(parsed.Second())*time.Second, nil
}
