package main

import (
	"NHL_Project/nhlApi"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	now := time.Now()

	rosterFile, err := os.OpenFile("rosters.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Error: Unable to open file 'rosters.txt': %v", err)
	}
	defer rosterFile.Close()

	wrt := io.MultiWriter(os.Stdout, rosterFile)

	log.SetOutput(wrt)

	teams, err := nhlApi.GetAllTeams()
	if err != nil {
		log.Fatalf("Error: Unable to fetch teams from API: %v", err)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(teams))

	results := make(chan []nhlApi.Player)

	for _, team := range teams {
		go func(team nhlApi.Team) {
			roster, err := nhlApi.GetRoster(team.ID)
			if err != nil {
				log.Fatalf("Error fetching roster %v: %v", team.Name, err)
			}

			results <- roster
			waitGroup.Done()
		}(team)
	}

	go func() {
		waitGroup.Wait()
		close(results)
	}()

	display(results)

	log.Printf("Application took %s to run.", time.Now().Sub(now).String())
}

func display(results chan []nhlApi.Player) {
	for roster := range results {
		for _, player := range roster {
			log.Println("------------------------------------------")
			log.Printf("Name: %s\n", player.Person.FullName)
			log.Printf("Position: %v\n", player.Position.Abbreviation)
			log.Printf("Jersey Number: %v\n", player.JerseyNumber)
			log.Println("------------------------------------------")
		}
	}
}
