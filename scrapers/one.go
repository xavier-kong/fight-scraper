package scrapers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/icza/gox/timex"
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

func getEventInfo(url string) types.Event {
	c := colly.NewCollector(
		colly.AllowedDomains("www.onefc.com"),
	)

	event := types.Event{
		Url: url,
		Org: "onefc",
	}

	c.OnHTML(".info-content", func(e *colly.HTMLElement) {
		event.Headline = e.ChildText(".title")
		fmt.Println(event.Headline)

		e.ForEach(".event-date-time", func(i int, h *colly.HTMLElement) {
			if (h.ChildText(".timezone") == "ICT") {
				dateString := createDateString(h.ChildText(".day"))
				timeString := createTimeString(h.ChildText(".time"))
				event.TimestampSeconds = createTimestamp(dateString, timeString)
			}
		})
	})

	c.Visit(url)

	return event
}

func fetchOneEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.onefc.com"),
	)

	var newEvents []types.Event
	var eventsToUpdate []types.Event

	eventUrls := fetchEventUrls(c)

	for _, url := range(eventUrls) {
		event := getEventInfo(url)

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

func createDateString(day string) string {
	dayOfWeekRegex := regexp.MustCompile(`\s\([^()]*\)`)
	dayMonthString := dayOfWeekRegex.ReplaceAllString(day, "")

	monthString := regexp.MustCompile(`^[A-Za-z]+`).FindString(dayMonthString)
	monthObj, _ := timex.ParseMonth(monthString)
	monthInt := int(monthObj)

	dayRegex := regexp.MustCompile(`[0-9]+`)
	dayString := dayRegex.FindString(dayMonthString)

	yearInt := time.Now().Year()

	dateString := fmt.Sprintf("%d-%02d-%02s", yearInt, monthInt, dayString)

	return dateString
}

func createTimeString(time string) string {
	ending := regexp.MustCompile(`AM|PM`).FindString(time)
	timeString := regexp.MustCompile(`AM|PM`).ReplaceAllString(time, "")

	if len(timeString) == 1 {
		timeString += ":00"
	}

	vals := strings.Split(timeString, ":")
	hourString, minString := vals[0], vals[1]

	hourInt, _ := strconv.Atoi(hourString)
	hourString = fmt.Sprintf("%02d", hourInt)

	timeString = hourString + ":" + minString + ending

	return timeString
}

func createTimestamp(day string, time string) int {
	fmt.Println(day, time)

	return 0
}
