package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/xavier-kong/fight-scraper/types"
	"github.com/xavier-kong/fight-scraper/scrapers"
)

var Database *gorm.DB

func main() {
	loadEnv()
	createDbClient()

	existingEvents := createExistingEventsMap(Database)
	newEvents, eventsToUpdate := scrapers.FetchNewEvents(existingEvents)
	go writeNewEventsToDb(Database, newEvents)
	go updateExistingEvents(Database, eventsToUpdate)
}

func handleError(err error) {
	log.Fatal(err)
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		handleError(errors.New("error loading .env file"))
	}
}

func createDbClient() {
	var err error
	//dsn := fmt.Sprintf("%s&parseTime=True", os.Getenv("DSN"))

	dsn := os.Getenv("DSN")

	Database, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{DisableForeignKeyConstraintWhenMigrating: true},
	)

	if err == nil {
		fmt.Println("Successfully connected to PlanetScale!")
	} else {
		handleError(err)
	}
}

func createExistingEventsMap(db *gorm.DB) map[string]map[string]types.Event {
	m := make(map[string]map[string]types.Event)
	var events []types.Event

	todaySecs := int(time.Now().UnixMilli() / 1000)

	result := db.Where("timestamp_seconds > ?", todaySecs).Find(&events)

	if result.Error != nil {
		handleError(result.Error)
	}

	for _, event := range events {
		if _, exists := m[event.Org]; !exists {
			m[event.Org] = make(map[string]types.Event)
		}

		m[event.Org][event.Name] = event
	}

	return m
}

func writeNewEventsToDb(db *gorm.DB, events []types.Event) {
	if (len(events) == 0) {
		fmt.Println("no new events...returning")
		return
	}

	result := db.Create(&events)

	if result.Error != nil {
		handleError(result.Error)
	}
}
