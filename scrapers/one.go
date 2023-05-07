package scrapers

import (
	"fmt"
	"strconv"
	"time"
	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
)

func fetchOneEvents(existingEvents map[string]bool) []types.Event {
	c := colly.NewCollector(
		colly.AllowedDomains("www.ufc.com"),
	)

	todaySecs := int(time.Now().UnixMilli() / 1000)

	var events []types.Event

	c.OnHTML(".c-card-event--result__info", func(e *colly.HTMLElement) {
		eventHeadline := e.ChildText(".c-card-event--result__headline")

		timestampString := e.ChildAttr(".c-card-event--result__date", "data-main-card-timestamp")
		timestampMs, err := strconv.Atoi(timestampString);

		if err != nil { fmt.Printf("error converting %s to int", e.ChildAttr(".c-card-event--result__date", "data-main-card-timestamp")); return
	}

	if timestampMs < todaySecs {
		return
	}

	eventUrlPath := e.ChildAttr("a", "href")

	eventUrl := "https://www.ufc.com" + eventUrlPath

	eventName := convertUrlToEventName(eventUrlPath)

	if _, exists := existingEvents[eventName]; exists {
		fmt.Printf("event: %s already exists", eventName)
		return
	}

	event := types.Event{
		Name: eventName,
		Headline: eventHeadline,
		TimestampSeconds: timestampMs,
		Url: eventUrl,
		Org: "ufc",
	}

	events = append(events, event)
})

c.Visit("https://www.ufc.com/events#events-list-upcoming")

return events
}
