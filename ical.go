package main

import (
	"log"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	uuid "github.com/gofrs/uuid"
)

type Calendar struct {
	cal *ics.Calendar
}

func ParseCalendar(text string) (*Calendar, error) {
	cal, err := ics.ParseCalendar(strings.NewReader(text))
	if err != nil {
		return nil, err
	}
	log.Printf("Parsed calendar")
	return &Calendar{cal: cal}, nil
}

func (c *Calendar) RemoveEventsByName(eventsToRemove []string) (*Calendar, error) {
	newCalendar := ics.NewCalendar()
	newCalendar.SetMethod(ics.MethodPublish)
	newCalendar.SetProductId("-//Pascal Riesinger//lecture-mod")
	for _, event := range c.cal.Events() {
		summary := event.GetProperty(ics.ComponentPropertySummary).Value
		startDate := event.GetProperty(ics.ComponentPropertyDtStart).Value
		endDate := event.GetProperty(ics.ComponentPropertyDtEnd).Value
		location := event.GetProperty(ics.ComponentPropertyLocation).Value
		// lastModified := event.GetProperty(ics.ComponentPropertyLastModified).Value
		organizer := event.GetProperty(ics.ComponentPropertyOrganizer).Value
		log.Printf("Found event %s (id: %s), starting at %s", summary, event.Id(), startDate)
		// FIXME: Find a non-quadratic way of removing things
		shouldBeRemoved := false
		for _, remove := range eventsToRemove {
			// FIXME: Pull out lowecasing of the 'removing' patterns of the loop
			if strings.Contains(strings.ToLower(summary), strings.ToLower(remove)) {
				shouldBeRemoved = true
				break
			}
		}
		if shouldBeRemoved {
			continue
		}
		newEvent := newCalendar.AddEvent(uuid.Must(uuid.NewV4()).String())
		newEvent.SetSummary(summary)
		startTime, err := parseICSDate(startDate)
		if err != nil {
			log.Printf("Could not parse start date: %s", startDate)
			return nil, err
		}
		newEvent.SetStartAt(startTime)
		endTime, err := parseICSDate(endDate)
		if err != nil {
			log.Printf("Could not parse end date: %s", endDate)
			return nil, err
		}
		newEvent.SetEndAt(endTime)
		newEvent.SetLocation(location)
		// TODO: Parse the time
		// newEvent.SetModifiedAt(lastModified)
		newEvent.SetOrganizer(organizer)
	}
	return &Calendar{cal: newCalendar}, nil
}

func (c *Calendar) Serialize() string {
	return c.cal.Serialize()
}

func parseICSDate(date string) (t time.Time, err error) {
	t, err = time.Parse("20060102T150405Z", date)
	if err != nil {
		// One more try, because the Z is dropped some times :facepalm:
		t, err = time.Parse("20060102T150405", date)
	}
	return
}
