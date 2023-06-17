package scrapers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"github.com/araddon/dateparse"
	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
)

type Bell struct{}

var b Bell

func fetchBellatorEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.bellator.com"),
	)

	currDate := time.Now()

	var newEvents []types.Event
	var eventsToUpdate []types.Event
	todaySecs := int(currDate.UnixMilli() / 1000)

	c.OnHTML("html", func(page *colly.HTMLElement) {
		linkSet := make(map[string]bool)

		page.ForEach("a", func(i int, aTag *colly.HTMLElement) {
			link := aTag.Attr("href")
			if strings.Contains(link, "/event/") && !strings.Contains(link, "ticketmaster") {
				linkSet[link] = true
			}
		})

		for link := range linkSet {
			event := b.fetchEventData(fmt.Sprintf("https://www.bellator.com%s", link))

			if event.TimestampSeconds > 0 && event.TimestampSeconds < todaySecs {
				fmt.Println(event.Headline, "past")
				continue
			}

			existingEventData, exists := existingEvents[event.Name]

			if !exists {
				newEvents = append(newEvents, event)
				continue
			}

			if existingEventData.TimestampSeconds == 0 || existingEventData.Headline != event.Headline {
				event.ID = existingEventData.ID
				eventsToUpdate = append(eventsToUpdate, event)
			}
		}
	})

	c.Visit("https://www.bellator.com/event")

	return newEvents, eventsToUpdate
}

func (b Bell) fetchEventData(link string) types.Event {
	event := types.Event{
		Org: "bellator",
		Url: link,
	}

	c := colly.NewCollector(
		colly.AllowedDomains("www.bellator.com"),
	)

	c.OnHTML("html", func(page *colly.HTMLElement) {
		page.ForEach("h1", func(i int, h1 *colly.HTMLElement) {
			if strings.Contains(h1.Attr("class"), "Titlestyles") {
				event.Headline = h1.ChildText("span")
			}
		})

		page.ForEach("h2", func(i int, h2 *colly.HTMLElement) {
			if strings.Contains(h2.Attr("class"), "BellatorNumber") {
				event.Name = h2.Text
			}
		})

		page.ForEach("time", func(i int, timeContainer *colly.HTMLElement) {
			dateTime := timeContainer.Attr("datetime")

			if string(dateTime[0]) == "P" {
				// P3DT15H10.926033333333333M
				re := regexp.MustCompile(`P(?P<Days>\d*)DT(?P<Hours>\d*)H(?P<Minutes>\d*).\d*`)
				groups := re.FindStringSubmatch(dateTime)

				if len(groups) != 4 {
					fmt.Printf("error with date time %s\n", dateTime)
					event.TimestampSeconds = 0
				} else {
					dayString, hoursString, minutesString := groups[1], groups[2], groups[3]

					daysInt, err := strconv.Atoi(dayString)
					hoursInt, err := strconv.Atoi(hoursString)
					minutesInt, err := strconv.Atoi(minutesString)
					if err != nil {
						handleError(err)
					}

					ts := time.Now().AddDate(0, 0, daysInt).
						Add(time.Hour*time.Duration(hoursInt) + time.Minute*time.Duration(minutesInt)).
						Round(time.Minute * 30)

					event.TimestampSeconds = int(ts.UnixMilli()) / 1000
				}
			}
		})
	})

	c.Visit(link)

	return event
}
