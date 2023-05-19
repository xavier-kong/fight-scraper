package scrapers

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
)

type Pfl struct {}

var pfl Pfl

func fetchPflEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.pflmma.com"),
	)

	var newEvents []types.Event
	var eventsToUpdate []types.Event

	c.OnHTML(".row", func(e *colly.HTMLElement) {



	})

	c.Visit(fmt.Sprintf("https://www.pflmma.com/season/%d", time.Now().Year()))

	return newEvents, eventsToUpdate
}

