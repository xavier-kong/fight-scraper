package scrapers

import (
	//"fmt"
	"fmt"
	"strings"
	"time"

	//"time"

	"github.com/araddon/dateparse"
	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
)

type Bell struct {}

func fetchBellatorEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.espn.com"),
	)

	currDate := time.Now()

	var newEvents []types.Event
	var eventsToUpdate []types.Event
	//var a Bell todaySecs := int(time.Now().UnixMilli() / 1000)

	c.OnHTML(".page-container", func(container *colly.HTMLElement) {
		container.ForEachWithBreak("tbody.Table__TBODY", func(i int, table *colly.HTMLElement) bool {
			table.ForEach("tr", func(j int, row *colly.HTMLElement) {
				event := types.Event {
					Headline: row.ChildText("td.event__col"),
					Name: strings.Split(row.ChildText("td.event__col"), ":")[0],
					Url: fmt.Sprintf("www.espn.com%s", row.ChildAttr("a", "href")),
					Org: "bellator",
				}

				dateString := fmt.Sprintf("%s %d", row.ChildText("td:nth-child(1)"), currDate.Year())
				timeString := fmt.Sprintf("%s GMT+0800", row.ChildText("td:nth-child(2)"))

				ts, err := dateparse.ParseAny(fmt.Sprintf("%s %s", dateString, timeString))

				if err != nil {
					fmt.Println(err)
				}

				event.TimestampSeconds = int(ts.UnixMilli()) / 1000

				if event.TimestampSeconds < todaySecs {
					fmt.Println(event.Headline, " past")
					return
				}

				existingEventData, exists := existingEvents[event.Name]

				if !exists {
					newEvents = append(newEvents, event)
					return
				}

				if (existingEventData.TimestampSeconds != event.TimestampSeconds ||
				existingEventData.Headline != event.Headline) {
					event.ID =  existingEventData.ID
					eventsToUpdate = append(eventsToUpdate, event)
				}
			})

			return false
		})
	})

	c.Visit("https://www.espn.com/mma/schedule/_/league/bellator")

	return newEvents, eventsToUpdate
}


