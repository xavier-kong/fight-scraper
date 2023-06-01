package scrapers

import (
	//"fmt"
	"time"
	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
)

type Bell struct {}

func fetchBellatorEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.espn.com"),
	)

	var newEvents []types.Event
	var eventsToUpdate []types.Event
	var a Bell

	todaySecs := int(time.Now().UnixMilli() / 1000)

	c.Visit("https://www.espn.com/mma/schedule/_/league/bellator")

	return newEvents, eventsToUpdate
}


