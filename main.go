package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	//загружаем конфиг
	cfg, err := LoadConfig("sunny_5_skiers/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	//открываем файл событий
	eventsFile, err := os.Open("sunny_5_skiers/events")
	if err != nil {
		log.Fatalf("Error opening events file: %v", err)
	}
	defer eventsFile.Close()

	//парсинг событий
	scanner := bufio.NewScanner(eventsFile)
	var events []Event
	for scanner.Scan() {
		line := scanner.Text()
		evt, err := ParseEvent(line)
		if err != nil {
			log.Printf("Skipping invalid event line: %s (error: %v)", line, err)
			continue
		}
		events = append(events, evt)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading events file: %v", err)
	}

	//обработка событий
	system := NewCompetitionSystem(cfg)
	for _, evt := range events {
		system.HandleEvent(evt)
	}

	//вывод результатов
	fmt.Println("\nEvent Log")
	for _, logEntry := range system.EventLog {
		fmt.Println(logEntry)
	}

	fmt.Println("\nFinal Report")
	system.PrintFinalReport()
}
