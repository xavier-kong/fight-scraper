package scrapers

import (
	"fmt"
	//"strconv"
	//"time"
	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
)

func fetchEventUrls(c *colly.Collector) []string {
	eventUrls := make([]string, 0)

	c.OnHTML("#upcoming-events-section", func(e *colly.HTMLElement) {
		e.ForEach(".simple-post-card", func(i int, h *colly.HTMLElement) {
			eventUrls = append(eventUrls, h.ChildAttr("a", "href"))
		})
	})

	c.Visit("https://www.onefc.com/events")

	return eventUrls;
}

func getEventInfo(c *colly.Collector, url string) types.Event {

}

func fetchOneEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.onefc.com"),
	)

	var newEvents []types.Event
	var eventsToUpdate []types.Event

	eventUrls := fetchEventUrls(c)

	for _, url := range(eventUrls) {
		event := getEventInfo(c, url)

		existingEventData, exists := existingEvents[event.Name]

		if exists {
			if (existingEventData.TimestampSeconds != event.TimestampSeconds ||
			existingEventData.Headline != event.Headline) {
				event.ID =  existingEventData.ID
				eventsToUpdate = append(eventsToUpdate, event)
			}
		} else {
			newEvents = append(newEvents, event)
		}
	}

	return newEvents, eventsToUpdate
}
