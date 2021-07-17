package main

import (
	"context"
	"log"
	"os"
	"path"

	"golang.org/x/sync/errgroup"
)

func main() {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = path.Join(os.Getenv("HOME"), ".config")
	}
	configFile := path.Join(configHome, "i3statuspp", "config.json")
	cfg, err := LoadConfig(configFile)
	if err != nil {
		log.Printf("[WARN] Unable to load config file: %s", err)
	} else {
		log.Println("[INFO] Successfully loaded config file:", configFile)
	}

	eventChan := make(chan ClickEvent, 32)

	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		var i3statuscmd string
		if cfg.General.I3StatusCMD != "" {
			i3statuscmd = cfg.General.I3StatusCMD
		} else {
			i3statuscmd = "i3status"
		}
		return WrapI3Status(i3statuscmd)
	})
	g.Go(func() error {
		return ReadEvents(os.Stdin, eventChan)
	})
	g.Go(func() error {
		return ProcessEvents(cfg, eventChan)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
