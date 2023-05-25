package scrapers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Pfl struct {}

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

			if len(parts) != 3 { // past return
				return
			}

			timestamp := pfl.getTimestamp(parts[0],  strings.Replace(parts[2], "ESPN ", "", -1))

			if timestamp < todaySecs { // past
				return
			}

			eventNameString := strings.ToLower(h.ChildText("p.text-red.font-weight-bold.m-0"))
			eventName := cases.Title(language.English, cases.NoLower).String(eventNameString)
			eventUrl := h.ChildAttr("a", "href")

			event := types.Event {
				TimestampSeconds: timestamp,
				Name: eventName,
				Headline: eventName,
				Url: eventUrl,
				Org: "pfl",
			}

			fmt.Println(event)

		})
	})

	c.Visit(fmt.Sprintf("https://www.pflmma.com/season/%d", time.Now().Year()))

	return newEvents, eventsToUpdate
}

func (p Pfl) getTimestamp(date string, time string) int {
	date = regexp.MustCompile(`Monday |Tuesday |Wednesday |Thursday |Friday |Saturday |Sunday `).ReplaceAllString(date, "")
	ts, err := dateparse.ParseAny(fmt.Sprintf("%s %s", date, strings.Replace(time, "PM", ":00PM", -1)))

	if err != nil {
		handleError(err)
	}

	return int(ts.UnixMilli()) / 1000
}
