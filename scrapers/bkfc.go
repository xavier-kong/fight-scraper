package scrapers

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/araddon/dateparse"
	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
)

type Bkfc struct {}

var bkfc Bkfc


func fetchBkfcEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	var newEvents []types.Event
	var eventsToUpdate []types.Event

	eventTimestamps := getEventTimestamps()

	fmt.Println(eventTimestamps)

	c := colly.NewCollector(
		colly.AllowedDomains("www.bkfc.com"),
	)

	c.OnHTML(".row.card-module", func(upcomingContainer *colly.HTMLElement) {
		upcomingContainer.ForEach("div.col-12.col-lg-4.mb-3", func(i int, eventCard *colly.HTMLElement) {
			event := types.Event { Org: "bfkc" }

			eventCard.ForEach(".card-text-events", func(i int, eventHeader *colly.HTMLElement) {
				eventNameString := strings.ToLower(eventHeader.ChildAttr("a", "title"))
				eventNameStringCased := cases.Title(language.English, cases.NoLower).String(eventNameString)
				event.Name = strings.Replace(eventNameStringCased, "Bkfc", "BKFC", -1)
				event.Headline = event.Name

				event.Url = fmt.Sprintf("https://www.bkfc.com%s", eventHeader.ChildAttr("a", "href"))



			})

			fmt.Println(event)
		})
	})

	c.Visit("https://www.bkfc.com/events")

	return newEvents, eventsToUpdate
}

func getEventTimestamps() map[string]int {
	c := colly.NewCollector(
		colly.AllowedDomains("www.itnwwe.com"),
	)

	tsMap := make(map[string]int)

	c.OnHTML("tbody", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(i int, row *colly.HTMLElement) {
			timeString := row.ChildText("td:nth-child(4)")

			if timeString == "" {
				return
			}

			eventNumber := regexp.MustCompile(`\d+`).FindString(row.ChildText("td:nth-child(2)"))
			dateTimeString := fmt.Sprintf("%s %s", row.ChildText("td:nth-child(1)"), strings.ReplaceAll(timeString, "8:00", "08:00"))

			ts, err := dateparse.ParseAny(dateTimeString)

			if err != nil {
				handleError(err)
			}


			fmt.Println("date", ts, eventNumber)
		})
	})

	c.Visit("https://www.itnwwe.com/mma/bkfc-events-schedule/")

	return tsMap
}
