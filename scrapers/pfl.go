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

	c.OnHTML(".container", func(e *colly.HTMLElement) {
		e.ForEach(".row.py-4", func(i int, h *colly.HTMLElement) {
			h.ForEach("p.font-oswald.font-weight-bold.m-0", func(i int, j *colly.HTMLElement) {
				fmt.Println(j.Text)
			})
		})
	})

	c.Visit(fmt.Sprintf("https://www.pflmma.com/season/%d", time.Now().Year()))

	return newEvents, eventsToUpdate
}

