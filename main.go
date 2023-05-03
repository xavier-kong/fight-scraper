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
	name string
	timestamp int
	headline string
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.ufc.com"),
	)

	today := time.Now()

	c.OnHTML(".c-card-event--result__info", func(e *colly.HTMLElement) {
		eventHeadline := e.ChildText(".c-card-event--result__headline")

		timestampString := e.ChildAttr(".c-card-event--result__date", "data-main-card-timestamp")
		timestamp, err := strconv.Atoi(timestampString);

		if err != nil {
			fmt.Printf("error converting %s to int", e.ChildAttr(".c-card-event--result__date", "data-main-card-timestamp"));
			return
		}

		eventTime := time.UnixMilli(int64(timestamp))

		if eventTime.Before(today) {
			fmt.Println(eventTime)
			fmt.Printf("event %s has past", eventHeadline)
			return
		}

		eventName := convertUrlToEventName(e.ChildAttr("a", "href"))

		event := Event{
			name: eventName,
			headline: eventHeadline,
			timestamp: timestamp,
		}

		fmt.Println(event)
	})

	c.Visit("https://www.ufc.com/events#events-list-upcoming")
}

func convertUrlToEventName(url string) string {
	res := strings.ReplaceAll(url, "/event/ufc", "UFC")
	res = strings.ReplaceAll(res, "-", " ")
	c := cases.Title(language.English, cases.NoLower)
	return c.String(res)
}
