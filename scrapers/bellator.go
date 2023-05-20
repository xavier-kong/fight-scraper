package scrapers

import (
	//"fmt"
	//"time"
	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
)

type Org struct {}

func (a Org) fetchEventUrls(c *colly.Collector) []string {
	eventUrls := make([]string, 0)

	c.OnHTML("#upcoming-events-section", func(e *colly.HTMLElement) {
		e.ForEach(".simple-post-card", func(i int, h *colly.HTMLElement) {
			eventUrls = append(eventUrls, h.ChildAttr("a", "href"))
		})
	})

	c.Visit("https://www.onefc.com/events")

	return eventUrls;
}

func fetchBellatorEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	/*c := colly.NewCollector(
		colly.AllowedDomains("www.bellator.com"),
	)*/

	var newEvents []types.Event
	var eventsToUpdate []types.Event

	//var a Org

	//eventUrls := a.fetchEventUrls(c)

	//todaySecs := int(time.Now().UnixMilli() / 1000)

	/*for _, url := range(eventUrls) {
		event := getEventInfo(url)

		if event.TimestampSeconds < todaySecs {
			fmt.Println(event.Headline, "past")
			continue
		}

		existingEventData, exists := existingEvents[event.Name]

		if !exists {
			newEvents = append(newEvents, event)
			continue
		}

		if (existingEventData.TimestampSeconds != event.TimestampSeconds ||
		existingEventData.Headline != event.Headline) {
			event.ID =  existingEventData.ID
			eventsToUpdate = append(eventsToUpdate, event)
		}
	}*/

	return newEvents, eventsToUpdate
}


