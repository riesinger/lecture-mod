package main

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	currentCalendar *Calendar

	mtx sync.RWMutex
)

// TODO: Make this dynamic
const (
	icalURL = "https://rapla.dhbw-karlsruhe.de/rapla?page=ical&user=freudenmann&file=TINF17B1"
)

func main() {
	go refresher()

	startRouter()
}

func refresher() {
	refresh()
	ticker := time.NewTicker(1 * time.Minute)

	log.Println("Starting refresh timer")

	for {
		select {
		case <-ticker.C:
			refresh()
		}
	}
}

func refresh() {
	log.Println("Refreshing...")

	resp, err := http.Get(icalURL)
	if err != nil {
		log.Fatalf("Error while refreshing: %v", err)
		return
	}

	defer resp.Body.Close()
	buff := new(bytes.Buffer)
	buff.ReadFrom(resp.Body)
	calendar, err := ParseCalendar(buff.String())
	if err != nil {
		log.Printf("Error parsing the calendar: %v", err)
		return
	}
	mtx.Lock()
	currentCalendar = calendar
	mtx.Unlock()

	log.Println("Calendar was refreshed")
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func generateCalendar(params map[string]string) ([]byte, error) {

	tagsString, ok := params["tags"]
	if !ok {
		log.Println("Got no tags, returning whole calendar")
		mtx.RLock()
		defer mtx.RUnlock()
		return []byte(currentCalendar.Serialize()), nil
	}

	tags := strings.Split(tagsString, "+")

	mtx.RLock()
	cal, err := currentCalendar.RemoveEventsByName(tags)
	mtx.RUnlock()
	if err != nil {
		log.Printf("Failed to remove events: %v", err)
		return nil, err
	}
	return []byte(cal.Serialize()), nil

}
