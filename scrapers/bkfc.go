package scrapers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/araddon/dateparse"
	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
)

type Bkfc struct{}

var bkfc Bkfc

func fetchBkfcEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	var newEvents []types.Event
	var eventsToUpdate []types.Event

	eventTimestamps := bkfc.getEventTimestamps()
	todaySecs := int(time.Now().UnixMilli() / 1000)

	c := colly.NewCollector(
		colly.AllowedDomains("www.bkfc.com"),
	)

	c.OnHTML(".row.card-module", func(upcomingContainer *colly.HTMLElement) {
		upcomingContainer.ForEach("div.col-12.col-lg-4.mb-3", func(i int, eventCard *colly.HTMLElement) {
			event := types.Event{Org: "bkfc"}

			eventCard.ForEach(".card-text-events", func(i int, eventHeader *colly.HTMLElement) {
				eventHeadlineString := strings.ToLower(eventHeader.ChildAttr("a", "title"))
				eventHeadlineStringCased := cases.Title(language.English, cases.NoLower).String(eventHeadlineString)
				event.Headline = strings.Replace(eventHeadlineStringCased, "Bkfc", "BKFC", -1)

				eventNumber := regexp.MustCompile(`\d+`).FindString(event.Headline)
				event.Name = fmt.Sprintf("BKFC %s", eventNumber)

				event.Url = fmt.Sprintf("https://www.bkfc.com%s", eventHeader.ChildAttr("a", "href"))

				eventTs, exist := eventTimestamps[eventNumber]

				if !exist { // use future incorrect timestamp that will be updated when scraper runs again
					fmt.Println("no ts for ", eventNumber)
					eventTs = int(time.Now().AddDate(0, 0, 7).UnixMilli() / 1000)
				}

				event.TimestampSeconds = eventTs
			})

			if event.TimestampSeconds < todaySecs {
				fmt.Println(event.Headline, " past")
				return
			}

			existingEventData, exists := existingEvents[event.Name]

			if !exists {
				newEvents = append(newEvents, event)
				return
			}

			if existingEventData.TimestampSeconds != event.TimestampSeconds ||
				existingEventData.Headline != event.Headline {
				event.ID = existingEventData.ID
				eventsToUpdate = append(eventsToUpdate, event)
			}
		})
	})

	c.Visit("https://www.bkfc.com/events")

	return newEvents, eventsToUpdate
}

func (b Bkfc) getEventTimestamps() map[string]int {
	c := colly.NewCollector(
		colly.AllowedDomains("www.itnwwe.com"),
	)

	tsMap := make(map[string]int)

	c.OnHTML("tbody", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(i int, row *colly.HTMLElement) {
			timeString := strings.ReplaceAll(row.ChildText("td:nth-child(4)"), "EST", "GMT-0400")

			if timeString == "" {
				return
			}

			eventNumber := regexp.MustCompile(`\d+`).FindString(row.ChildText("td:nth-child(2)"))
			dateStringParts := strings.Split(row.ChildText("td:nth-child(1)"), " ")
			day, monthStr, year := dateStringParts[0], dateStringParts[1], dateStringParts[2]
			monthInt := bkfc.monthNameToInt(strings.ToLower(monthStr))
			dateTimeString := fmt.Sprintf("%s-%02d-%s %s", year, monthInt, day, timeString)

			ts, err := dateparse.ParseAny(dateTimeString)
			if err != nil {
				fmt.Println("error parsing", dateTimeString)
				tsMap[eventNumber] = int(time.Now().AddDate(0, 0, 7).UnixMilli()) / 1000
				return
			}

			tsMap[eventNumber] = int(ts.UnixMilli()) / 1000
		})
	})

	c.Visit("https://www.itnwwe.com/mma/bkfc-events-schedule/")

	return tsMap
}

func (b Bkfc) monthNameToInt(month string) int {
	nameInt := map[string]int{
		"january":   1,
		"february":  2,
		"march":     3,
		"april":     4,
		"may":       5,
		"june":      6,
		"july":      7,
		"august":    8,
		"september": 9,
		"october":   10,
		"november":  11,
		"december":  12,
	}

	val, exists := nameInt[month]

	if !exists {
		val = 1
	}

	return val
}
