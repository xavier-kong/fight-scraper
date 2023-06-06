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
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"context"
)

var Database *gorm.DB

func verifyOrigin(req *events.LambdaFunctionURLRequest) bool {
	isVerified := false;

	token, ok := req.Headers["my-precious-token"]

	if !ok || len(token) == 0 {
		fmt.Println("Error: no token found in header")
		return false
	}

	secret := os.Getenv("FIGHT_SCRAPER_SECRET")

	if secret == "" {
		fmt.Println("no secret found")
		return false
	}

	isVerified = token == secret

	if !isVerified {
		fmt.Println("hashes are not equal")
	} else {
		fmt.Println("signature has been verified")
	}

	return isVerified
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

func createExistingEventsMap() map[string]map[string]types.Event {
	m := make(map[string]map[string]types.Event)
	var events []types.Event

	todaySecs := int(time.Now().UnixMilli() / 1000)

	result := Database.Where("timestamp_seconds > ?", todaySecs).Find(&events)

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

func writeNewEventsToDb(events []types.Event) {
	if (len(events) == 0) {
		fmt.Println("no new events...returning")
		return
	}

	result := Database.Create(&events)

	if result.Error != nil {
		handleError(result.Error)
	}
}

func updateExistingEvents(eventsToUpdate []types.Event) {
	for _, event := range eventsToUpdate {
		Database.Save(&event)
	}
}

func logScrape(numNewEvents int, numEventsToUpdate int) {
	log := types.Log {
		Type: fmt.Sprintf("found %d new events and updated %d events", numNewEvents, numEventsToUpdate),
		TimestampSeconds: int(time.Now().UnixMilli()) / 1000,
	}

	result := Database.Create(&log)

	if (result.Error != nil) {
		handleError(result.Error)
	} else {
		fmt.Println("logged at ", log.TimestampSeconds)
	}
}

func handleRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (string, error) {
	isVerified := verifyOrigin(&req)

	if isVerified == false {
		return "verification error", nil
	}

	loadEnv()
	createDbClient()

	existingEvents := createExistingEventsMap()
	newEvents, eventsToUpdate := scrapers.FetchNewEvents(existingEvents)
	writeNewEventsToDb(newEvents)
	updateExistingEvents(eventsToUpdate)
	logScrape(len(newEvents), len(eventsToUpdate))

	return "done", nil
}

func main() {
	lambda.Start(handleRequest)
}
