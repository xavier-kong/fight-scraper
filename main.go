package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/gocolly/colly"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Event struct {
	id int
	name string
	timestampSeconds int
	headline string
	url string
	org string
}

func main() {
	existingEvents := createExistingEventsMap()
	newEvents := make([]Event, 0)

	for _, event := range fetchEvents() {
		if _, exists := existingEvents[event.name]; !exists {
			newEvents = append(newEvents, event)
		}
	}

	writeNewEventsToDb(newEvents)
}

func convertUrlToEventName(url string) string {
	res := strings.ReplaceAll(url, "/event/ufc", "UFC")
	res = strings.ReplaceAll(res, "-", " ")
	c := cases.Title(language.English, cases.NoLower)
	return c.String(res)
}

func createExistingEventsMap() map[string]bool {
	m := make(map[string]bool)

	// fetch all future from db

	return m
}

func writeNewEventsToDb(events []Event) {

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
			fmt.Printf("event %s has past\n", eventHeadline)
			return
		}

		eventUrlPath := e.ChildAttr("a", "href")

		eventUrl := "https://www.ufc.com" + eventUrlPath

		eventName := convertUrlToEventName(eventUrlPath)

		event := Event{
			name: eventName,
			headline: eventHeadline,
			timestampSeconds: timestampMs,
			url: eventUrl,
			org: "UFC",
		}

		events = append(events, event)
	})

	c.Visit("https://www.ufc.com/events#events-list-upcoming")

	return events;
}
