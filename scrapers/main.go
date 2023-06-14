package scrapers

import (
	"log"

	"github.com/xavier-kong/fight-scraper/types"
)

func FetchNewEvents(existingEvents map[string]map[string]types.Event) (newEvents []types.Event, eventsToUpdate []types.Event) {
	ufcNewEvents, ufceEventsToUpdate := fetchUfcEvents(existingEvents["ufc"])
	newEvents, eventsToUpdate = append(newEvents, ufcNewEvents...), append(eventsToUpdate, ufceEventsToUpdate...)

	oneNewEvents, oneEventsToUpdate := fetchOneEvents(existingEvents["one"])
	newEvents, eventsToUpdate = append(newEvents, oneNewEvents...), append(eventsToUpdate, oneEventsToUpdate...)

	bellNewEvents, bellEventsToUpdate := fetchBellatorEvents(existingEvents["bellator"])
	newEvents, eventsToUpdate = append(newEvents, bellNewEvents...), append(eventsToUpdate, bellEventsToUpdate...)

	bkfcNewEvents, bkfcEventsToUpdate := fetchBkfcEvents(existingEvents["bkfc"])
	newEvents, eventsToUpdate = append(newEvents, bkfcNewEvents...), append(eventsToUpdate, bkfcEventsToUpdate...)

	pflNewEvents, pflEventsToUpdate := fetchPflEvents(existingEvents["pfl"])
	newEvents, eventsToUpdate = append(newEvents, pflNewEvents...), append(eventsToUpdate, pflEventsToUpdate...)

	return
}

func handleError(err error) {
	log.Fatal(err)
}
