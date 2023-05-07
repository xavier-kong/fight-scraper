package scrapers

import ("github.com/xavier-kong/fight-scraper/types")

func FetchNewEvents(existingEvents map[string]map[string]types.Event) (newEvents []types.Event, eventsToUpdate []types.Event) {

	ufcNewEvents, ufcEventsToUpdate := fetchUfcEvents(existingEvents["ufc"])
	newEvents = append(newEvents , ufcNewEvents...)
	eventsToUpdate = append(eventsToUpdate, ufcEventsToUpdate...)

	return;
}
