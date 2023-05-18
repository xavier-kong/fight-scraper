package scrapers

import (
	"github.com/xavier-kong/fight-scraper/types"
	"github.com/gocolly/colly"
)

func fetchPflEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.pflmma.com"),
	)

	var newEvents []types.Event
	var eventsToUpdate []types.Event



	return newEvents, eventsToUpdate
}
