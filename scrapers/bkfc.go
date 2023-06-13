package scrapers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/xavier-kong/fight-scraper/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Bkfc struct{}

var bkfc Bkfc

func fetchBkfcEvents(existingEvents map[string]types.Event) ([]types.Event, []types.Event) {
	var newEvents []types.Event
	var eventsToUpdate []types.Event

	todaySecs := int(time.Now().UnixMilli() / 1000)

	c := colly.NewCollector(
		colly.AllowedDomains("www.bkfc.com"),
	)

	c.OnHTML(".row.card-module", func(upcomingContainer *colly.HTMLElement) {
		upcomingContainer.ForEach("div.col-12.col-lg-4.mb-3", func(i int, eventCard *colly.HTMLElement) {
			event := types.Event{Org: "bkfc"}

			eventCard.ForEach(".card-text-events", func(i int, eventHeader *colly.HTMLElement) {
				eventHeadlineString := strings.ToLower(eventHeader.ChildAttr("a", "title"))
				eventHeadlineStringCased := cases.Title(language.English, cases.NoLower).String(eventHeadlineString)
				event.Headline = strings.Replace(eventHeadlineStringCased, "Bkfc", "BKFC", -1)

				eventNumber := regexp.MustCompile(`\d+`).FindString(event.Headline)
				event.Name = fmt.Sprintf("BKFC %s", eventNumber)

				event.Url = fmt.Sprintf("https://www.bkfc.com%s", eventHeader.ChildAttr("a", "href"))

				event.TimestampSeconds = bkfc.getEventTimestamp(event.Url)
			})

			if event.TimestampSeconds < todaySecs {
				fmt.Println(event.Headline, " past")
				return
			}

			existingEventData, exists := existingEvents[event.Name]

			if !exists {
				newEvents = append(newEvents, event)
				return
			}

			if existingEventData.TimestampSeconds == 0 ||
				existingEventData.Headline != event.Headline {
				event.ID = existingEventData.ID
				eventsToUpdate = append(eventsToUpdate, event)
			}
		})
	})

	c.Visit("https://www.bkfc.com/events")

	return newEvents, eventsToUpdate
}

func (b Bkfc) getEventTimestamp(url string) int {
	dateTimeString := bkfc.getDateTimeString(url)

	if dateTimeString == "" {
		panic(fmt.Sprintf("url has failed %s", url))
	}

	parts := strings.Split(dateTimeString, "T")

	if len(parts) != 2 {
		panic(fmt.Sprintf("dtstring FAILED %s", dateTimeString))
	}

	dateString, timeStringWithOffset := parts[0], parts[1]

	timeString := strings.ReplaceAll(timeStringWithOffset, "+00:00", "")

	timeStamp, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", dateString, timeString))
	if err != nil {
		panic(fmt.Sprintf("FAILED %s %s", dateString, timeString))
	}

	return int(timeStamp.UnixMilli()) / 1000
}

func (b Bkfc) getDateTimeString(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// Convert the body to type string
	sb := string(body)

	scanner := bufio.NewScanner(strings.NewReader(sb))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "var date = moment(") {
			line = strings.ReplaceAll(line, "var date = moment(\"", "")
			line = strings.ReplaceAll(line, "\");", "")
			return strings.Trim(line, " ")
		}
	}

	return ""
}
