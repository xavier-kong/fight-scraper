package scrapers

import (
	"log"

	"github.com/xavier-kong/fight-scraper/types"
)

func FetchNewEvents(existingEvents map[string]map[string]types.Event) (newEvents []types.Event, eventsToUpdate []types.Event) {


	return;
}

func handleError(err error) {
	log.Fatal(err)
}
