package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	currentCalendar string

	mtx sync.RWMutex
)

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
	mtx.Lock()
	currentCalendar = buff.String()
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
		mtx.Lock()
		defer mtx.Unlock()
		return []byte(currentCalendar), nil
	}

	tags := strings.Split(tagsString, "+")

	builder := strings.Builder{}
	buf := strings.Builder{}

	buffering := false

	mtx.RLock()
	defer mtx.RUnlock()
	scanner := bufio.NewScanner(bytes.NewBufferString(currentCalendar))
	scanner.Split(ScanCRLF)

	numLines := 0

	for scanner.Scan() {
		numLines++
		l := scanner.Text()
		if buffering {
			buf.WriteString(l)
			buf.WriteString("\r\n")
			// Check if the current line terminates a VEVENT
			if strings.Contains(l, "END:VEVENT") {
				// If the currently buffered VEVENT contains any of the tags, dump it
				for _, v := range tags {
					p, err := url.PathUnescape(v)
					if err != nil {
						return nil, err
					}
					if strings.Contains(strings.ToLower(buf.String()), p) {
						buf.Reset()
						break
					}
				}

				if buf.String() != "" {
					builder.WriteString(buf.String())
				}

				buf.Reset()
				buffering = false
			}
		} else {
			// We currently are either in the metadata or something went very wrong!
			if strings.Contains(l, "BEGIN:VEVENT") {
				// Start buffering each event
				buffering = true
				buf.WriteString(l)
				buf.WriteString("\r\n")
			} else {
				// This should be metadata
				builder.WriteString(l)
				builder.WriteString("\r\n")
			}
		}
	}

	return []byte(builder.String()), nil

}
