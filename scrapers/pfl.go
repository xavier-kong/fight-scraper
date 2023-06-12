package scrapers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Pfl struct{}

var pfl Pfl

func fetchPflEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.pflmma.com"),
	)

	var newEvents []types.Event
	var eventsToUpdate []types.Event

	todaySecs := int(time.Now().UnixMilli() / 1000)

	c.OnHTML(".container", func(e *colly.HTMLElement) {
		e.ForEach(".row.py-4", func(i int, h *colly.HTMLElement) {
			dateTimeString := h.ChildText("p.font-oswald.font-weight-bold.m-0")
			parts := strings.Split(dateTimeString, " | ")

			if len(parts) <= 1 { // past return
				return
			}

			dateString, timeString := parts[0], ""

			if len(parts) == 3 {
				timeString = parts[2]
			}

			timestamp := pfl.getTimestamp(dateString, strings.Replace(timeString, "ESPN ", "", -1))

			if timestamp > 0 && timestamp < todaySecs { // past
				return
			}

			eventNameString := strings.ToLower(h.ChildText("p.text-red.font-weight-bold.m-0"))
			eventName := cases.Title(language.English, cases.NoLower).String(eventNameString)
			eventUrl := h.ChildAttr("a", "href")

			event := types.Event{
				TimestampSeconds: timestamp,
				Name:             eventName,
				Headline:         eventName,
				Url:              eventUrl,
				Org:              "pfl",
			}

			existingEventData, exists := existingEvents[event.Name]

			if !exists {
				newEvents = append(newEvents, event)
			} else if existingEventData.TimestampSeconds != event.TimestampSeconds ||
				existingEventData.Headline != event.Headline {
				event.ID = existingEventData.ID
				eventsToUpdate = append(eventsToUpdate, event)
			}
		})
	})

	c.Visit(fmt.Sprintf("https://www.pflmma.com/season/%d", time.Now().Year()))

	return newEvents, eventsToUpdate
}

func (p Pfl) getTimestamp(date string, timeString string) int {
	if timeString == "" {
		return 0
	}

	parts := strings.Split(date, " ")
	day := strings.ReplaceAll(parts[2], ",", "")
	month, year := parts[1], parts[3]

	if strings.Contains(timeString, "ET") {
		timeString = strings.ReplaceAll(timeString, "ET", "-0400")
	} else {
		fmt.Printf("unknown timezone detected in %s\n", timeString)
		return 0
	}

	ts, err := time.Parse("Jan 02 2006 15PM -0700", fmt.Sprintf("%s %s %s %s", month, day, year, timeString))
	if err != nil {
		handleError(err)
	}

	return int(ts.UnixMilli()) / 1000
}
