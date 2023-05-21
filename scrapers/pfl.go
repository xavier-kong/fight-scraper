package scrapers

import (
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
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
				parts := strings.Split(j.Text, " |")
				t, err := dateparse.ParseAny(parts[0])
				fmt.Println("pre", parts[0])

				if err != nil {
					fmt.Println("err", err)
				}

				fmt.Println("post", t.Year(), t.Month(), t.Day())

			})
		})
	})

	c.Visit(fmt.Sprintf("https://www.pflmma.com/season/%d", time.Now().Year()))

	return newEvents, eventsToUpdate
}

func (p Pfl) getTimestamp(date string) int {


}
