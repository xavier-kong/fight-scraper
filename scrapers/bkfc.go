package scrapers

import (
	"fmt"
	"strings"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
)

type Bkfc struct {}

var bkfc Bkfc


func fetchBkfcEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	var newEvents []types.Event
	var eventsToUpdate []types.Event

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

				event.TimestampSeconds = bkfc.fetchTimestamp(event.Url)
			})
		})
	})

	c.Visit("https://www.bkfc.com/events")


	return newEvents, eventsToUpdate
}

func (bkfc Bkfc) fetchTimestamp(url string) int {
	fmt.Println(url)
	ts := 0

	c := colly.NewCollector(
		colly.AllowedDomains("www.bkfc.com"),
	)

	c.OnHTML(".events-show-rails-date-im-timezone", func(dateHeader *colly.HTMLElement) {
		fmt.Println(dateHeader.Request)
	})

	c.Visit(url)

	return ts
}
