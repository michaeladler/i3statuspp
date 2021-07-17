package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"unicode"

	"github.com/google/shlex"
)

// ClickEvent as defined in https://i3wm.org/docs/i3bar-protocol.html
type ClickEvent struct {

	// Name of the block, if set
	Name string `json:"name"`

	// Instance of the block, if set
	Instance string `json:"instance"`

	// X11 button ID (for example 1 to 3 for left/middle/right mouse button)
	Button int `json:"button"`

	// An array of the modifiers active when the click occurred. The order in which modifiers are listed is not guaranteed.
	Modifiers []string `json:"modifiers"`

	// X11 root window coordinates where the click occurred
	X int `json:"x"`
	Y int `json:"y"`

	// Coordinates where the click occurred, with respect to the top left corner of the block
	RelativeX int `json:"relative_x"`
	RelativeY int `json:"relative_y"`

	// Coordinates relative to the current output where the click occurred
	OutputX int `json:"output_x,omitempty"`
	OutputY int `json:"output_y,omitempty"`

	//Width and height (in px) of the block
	Width  int `json:"width"`
	Height int `json:"height"`
}

func ReadEvents(rd io.Reader, outputCh chan ClickEvent) error {
	r := bufio.NewReader(rd)
	log.Println("[INFO] Waiting for events...")
	for {
	IgnoreChars:
		// Ignore unwanted chars first
		// idea borrowed from https://github.com/vincent-petithory/i3cat/blob/master/clickevents.go
		for {
			ruune, _, err := r.ReadRune()
			if err != nil {
				log.Println(err)
				break IgnoreChars
			}
			switch {
			case unicode.IsSpace(ruune):
				// Loop again
			case ruune == '[':
				// Loop again
			case ruune == ',':
				break IgnoreChars
			default:
				_ = r.UnreadRune()
				break IgnoreChars
			}
		}

		line, err := r.ReadBytes('\n')
		if err == io.EOF {
			log.Println("Reached EOF")
			parseClickEvent(outputCh, line)
			return io.EOF
		}
		if err != nil {
			log.Println("Error reading line. Skipping line.")
			continue
		}

		parseClickEvent(outputCh, line)
	}
}

type CacheKey struct {
	ID     string
	Button int
}

func ProcessEvents(cfg *Config, inputCh chan ClickEvent) error {
	cache := make(map[CacheKey][]string)

	for event := range inputCh {
		if cfg == nil {
			continue
		}
		for _, r := range cfg.Rules {
			if r.Name != "" && r.Name != event.Name {
				continue
			}
			if r.Instance != "" && r.Instance != event.Instance {
				continue
			}
			log.Printf("[DEBUG] Rule %s does match", r.ID)
			cacheKey := CacheKey{ID: r.ID, Button: event.Button}
			actionSplit, hit := cache[cacheKey]
			if !hit {
				log.Printf("[DEBUG] Cache miss for %v", cacheKey)
				var err error
				if action, found := r.Actions[strconv.Itoa(event.Button)]; found {
					actionSplit, err = shlex.Split(action)
					if err != nil {
						log.Printf("Error splitting action '%s': %s", action, err)
						continue
					}
					log.Printf("[DEBUG] Populating cache: key=%v, value=%v", cacheKey, actionSplit)
					cache[cacheKey] = actionSplit
				}
			}

			if len(actionSplit) > 0 {
				cmd := exec.Command(actionSplit[0], actionSplit[1:]...)
				cmd.Stdout = os.Stderr
				cmd.Stderr = os.Stderr
				go func() {
					if err := cmd.Run(); err != nil {
						log.Println("[ERROR]", err)
					}
				}()
			}
		}
	}
	return nil
}

func parseClickEvent(outputCh chan ClickEvent, line []byte) {
	var ce ClickEvent
	if err := json.Unmarshal(line, &ce); err != nil {
		log.Printf("Invalid JSON input: %v\n", err)
		return
	}
	log.Printf("Parsed click event %v\n", ce)
	outputCh <- ce
}
