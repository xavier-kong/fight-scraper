package scrapers

import (
	"fmt"
	"strconv"
	"time"
	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
	"strings"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func fetchUfcEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.ufc.com"),
	)

	todaySecs := int(time.Now().UnixMilli() / 1000)

	var newEvents []types.Event
	var eventsToUpdate []types.Event

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

		existingEventData, exists := existingEvents[eventName]

		event := types.Event{
			Name: eventName,
			Headline: eventHeadline,
			TimestampSeconds: timestampMs,
			Url: eventUrl,
			Org: "ufc",
		}

		if exists {
			if (existingEventData.TimestampSeconds != event.TimestampSeconds ||
			existingEventData.Headline != event.Headline) {
				event.ID =  existingEventData.ID
				eventsToUpdate = append(eventsToUpdate, event)
			}
		} else {
			newEvents = append(newEvents, event)
		}
	})

	c.Visit("https://www.ufc.com/events#events-list-upcoming")

	return newEvents, eventsToUpdate
}

func convertUrlToEventName(url string) string {
	res := strings.ReplaceAll(url, "/event/ufc", "UFC")
	res = strings.ReplaceAll(res, "-", " ")
	c := cases.Title(language.English, cases.NoLower)
	return c.String(res)
}