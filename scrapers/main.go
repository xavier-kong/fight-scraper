package scrapers

import ("github.com/xavier-kong/fight-scraper/types")

func FetchNewEvents(existingEvents map[string]map[string]bool) []types.Event {
	events := make([]types.Event, 0)

	ufcEvents := fetchUfcEvents(existingEvents["ufc"])
	events = append(events, ufcEvents...)

	return events;
}
