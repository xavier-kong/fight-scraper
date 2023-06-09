package scrapers

import (
	"fmt"
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
		colly.AllowedDomains("www.espn.com"),
	)

	currDate := time.Now()

	var newEvents []types.Event
	var eventsToUpdate []types.Event
	todaySecs := int(time.Now().UnixMilli() / 1000)

	c.OnHTML(".page-container", func(container *colly.HTMLElement) {
		container.ForEachWithBreak("tbody.Table__TBODY", func(i int, table *colly.HTMLElement) bool {
			table.ForEach("tr", func(j int, row *colly.HTMLElement) {
				event := types.Event{
					Headline: row.ChildText("td.event__col"),
					Name:     strings.Split(row.ChildText("td.event__col"), ":")[0],
					Url:      fmt.Sprintf("www.espn.com%s", row.ChildAttr("a", "href")),
					Org:      "bellator",
				}

				dateString := fmt.Sprintf("%s %d", row.ChildText("td:nth-child(1)"), currDate.Year())
				event.TimestampSeconds = b.convertToTimestamp(dateString, row.ChildText("td:nth-child(2)"))

				if event.TimestampSeconds < todaySecs {
					fmt.Println(event.Headline, " past")
					return
				}

				existingEventData, exists := existingEvents[event.Name]

				if !exists {
					newEvents = append(newEvents, event)
					return
				}

				if existingEventData.TimestampSeconds != event.TimestampSeconds ||
					existingEventData.Headline != event.Headline {
					event.ID = existingEventData.ID
					eventsToUpdate = append(eventsToUpdate, event)
				}
			})

			return false
		})
	})

	c.Visit("https://www.espn.com/mma/schedule/_/league/bellator")

	return newEvents, eventsToUpdate
}

func (b Bell) convertToTimestamp(date string, timeString string) int {
	fmt.Println("input", date, timeString)
	dateParts := strings.Split(date, " ")
	day, _ := strconv.Atoi(dateParts[1])

	timeParts := strings.Split(timeString, " ")
	hoursMinutes := strings.Split(timeParts[0], ":")

	hourString, minutes := hoursMinutes[0], hoursMinutes[1]

	hour, _ := strconv.Atoi(hourString)

	if timeParts[1] == "PM" && hour < 12 {
		hour += 12
	}

	dateTimeString := fmt.Sprintf("%02d %s %s %02d:%s GMT", day, dateParts[0], dateParts[2][2:], hour, minutes)
	ts, err := time.Parse("02 Jan 06 15:04 MST", dateTimeString)
	if err != nil {
		handleError(err)
	}

	return int(ts.UnixMilli() / 1000)
}
