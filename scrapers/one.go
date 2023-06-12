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

type One struct{}

var one One

func (a One) fetchEventUrls(c *colly.Collector) []string {
	eventUrls := make([]string, 0)

	c.OnHTML("#upcoming-events-section", func(e *colly.HTMLElement) {
		e.ForEach(".simple-post-card", func(i int, h *colly.HTMLElement) {
			eventUrls = append(eventUrls, h.ChildAttr("a", "href"))
		})
	})

	c.Visit("https://www.onefc.com/events")

	return eventUrls
}

func (a One) getEventInfo(url string) types.Event {
	c := colly.NewCollector(
		colly.AllowedDomains("www.onefc.com"),
	)

	event := types.Event{
		Url: url,
		Org: "one",
	}

	c.OnHTML(".info-content", func(e *colly.HTMLElement) {
		event.Headline = e.ChildText(".title")
		parts := strings.Split(event.Headline, ":")
		event.Name = strings.ReplaceAll(parts[0], "on Prime Video", "")

		e.ForEachWithBreak(".event-date-time", func(i int, h *colly.HTMLElement) bool {
			dateString := one.createDateString(h.ChildText(".day"))
			timeString := one.createTimeString(h.ChildText(".time"))
			timezoneString := h.ChildText(".timezone")
			offset := strings.Split(h.ChildAttr(".timezone.hint", "data-hint"), " GMT ")[0]

			t, err := time.Parse("2006-01-02 03:04PM -07:00", fmt.Sprintf("%s %s %s", dateString, timeString, offset))
			if err != nil { // use future incorrect timestamp that will be updated when scraper runs again
				fmt.Println("error parsing", dateString, timeString, timezoneString)
				event.TimestampSeconds = 0
				return false
			}

			event.TimestampSeconds = int(t.UnixMilli()) / 1000

			return false
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

	eventUrls := one.fetchEventUrls(c)

	todaySecs := int(time.Now().UnixMilli() / 1000)

	for _, url := range eventUrls {
		event := one.getEventInfo(url)

		if event.TimestampSeconds > 0 && event.TimestampSeconds < todaySecs {
			fmt.Println(event.Headline, "past")
			continue
		}

		existingEventData, exists := existingEvents[event.Name]

		if !exists {
			newEvents = append(newEvents, event)
			continue
		}

		if existingEventData.TimestampSeconds != event.TimestampSeconds ||
			existingEventData.Headline != event.Headline {
			event.ID = existingEventData.ID
			eventsToUpdate = append(eventsToUpdate, event)
		}
	}

	return newEvents, eventsToUpdate
}

func (a One) createDateString(day string) string {
	dayOfWeekRegex := regexp.MustCompile(`\s\([^()]*\)`)
	dayMonthString := dayOfWeekRegex.ReplaceAllString(day, "")

	monthString := regexp.MustCompile(`^[A-Za-z]+`).FindString(dayMonthString)
	monthObj, err := timex.ParseMonth(monthString)
	if err != nil {
		monthObj, _ = timex.ParseMonth("Jan")
	}

	monthInt := int(monthObj)

	dayRegex := regexp.MustCompile(`[0-9]+`)
	dayString := dayRegex.FindString(dayMonthString)

	yearInt := time.Now().Year()

	dateString := fmt.Sprintf("%d-%02d-%02s", yearInt, monthInt, dayString)

	return dateString
}

func (a One) createTimeString(time string) string {
	ending := regexp.MustCompile(`AM|PM`).FindString(time)
	timeString := regexp.MustCompile(`AM|PM`).ReplaceAllString(time, "")

	if len(timeString) == 1 {
		timeString += ":00"
	}

	vals := strings.Split(timeString, ":")
	hourString, minString := vals[0], vals[1]

	hourInt, err := strconv.Atoi(hourString)
	if err != nil {
		hourInt = 20
	}

	hourString = fmt.Sprintf("%02d", hourInt)

	timeString = hourString + ":" + minString + ending

	return timeString
}
