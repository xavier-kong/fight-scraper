package scrapers

import (
	"github.com/xavier-kong/fight-scraper/types"
	"log"
)

func FetchNewEvents(existingEvents map[string]map[string]types.Event) (newEvents []types.Event, eventsToUpdate []types.Event) {
	fetchBkfcEvents(existingEvents["bkfc"])

	return;
}

func handleError(err error) {
	log.Fatal(err)
}
