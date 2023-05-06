package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Event struct {
	ID uint
	Name string `gorm:"size: 255; not null;" json:"name"`
	TimestampSeconds int `gorm:"type: numeric; not null;" json:"timestamp_seconds"`
	Headline string `gorm:"size: 255; not null;" json:"headline"`
	Url string `gorm:"size: 255; not null;" json:"url"`
	Org string `gorm:"size: 255; not null;" json:"org"`
}

var Database *gorm.DB

func main() {
	loadEnv()
	createDbClient()

	existingEvents := createExistingEventsMap(Database)
	newEvents := filterOutOldEvents(existingEvents)
	writeNewEventsToDb(Database, newEvents)
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

func convertUrlToEventName(url string) string {
	res := strings.ReplaceAll(url, "/event/ufc", "UFC")
	res = strings.ReplaceAll(res, "-", " ")
	c := cases.Title(language.English, cases.NoLower)
	return c.String(res)
}

func createExistingEventsMap(db *gorm.DB) map[string]bool {
	m := make(map[string]bool)
	var events []Event

	todaySecs := int(time.Now().UnixMilli() / 1000)

	result := db.Where("timestamp_seconds > ?", todaySecs).Find(&events)

	if result.Error != nil {
		handleError(result.Error)
	}

	for _, event := range events {
		m[event.Name] = true
	}

	return m
}

func filterOutOldEvents(existingEvents map[string]bool) []Event {
	newEvents := make([]Event, 0)

	for _, event := range fetchEvents() {
		if _, exists := existingEvents[event.Name]; !exists {
			newEvents = append(newEvents, event)
		}
	}

	return newEvents
}

func writeNewEventsToDb(db *gorm.DB, events []Event) {
	if (len(events) == 0) {
		fmt.Println("no new events...returning")
		return
	}

	result := db.Create(&events)

	if result.Error != nil {
		handleError(result.Error)
	}
}

func fetchEvents() []Event {
	c := colly.NewCollector(
		colly.AllowedDomains("www.ufc.com"),
	)

	todaySecs := int(time.Now().UnixMilli() / 1000)

	events := make([]Event, 0)

	c.OnHTML(".c-card-event--result__info", func(e *colly.HTMLElement) {
		eventHeadline := e.ChildText(".c-card-event--result__headline")

		timestampString := e.ChildAttr(".c-card-event--result__date", "data-main-card-timestamp")
		timestampMs, err := strconv.Atoi(timestampString);

		if err != nil {
			fmt.Printf("error converting %s to int", e.ChildAttr(".c-card-event--result__date", "data-main-card-timestamp"));
			return
		}

		if timestampMs < todaySecs {
			return
		}

		eventUrlPath := e.ChildAttr("a", "href")

		eventUrl := "https://www.ufc.com" + eventUrlPath

		eventName := convertUrlToEventName(eventUrlPath)

		event := Event{
			Name: eventName,
			Headline: eventHeadline,
			TimestampSeconds: timestampMs,
			Url: eventUrl,
			Org: "UFC",
		}

		events = append(events, event)
	})

	c.Visit("https://www.ufc.com/events#events-list-upcoming")

	return events;
}
