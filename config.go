package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	General General `json:"general"`
	Rules   []Rule  `json:"rules"`
}

type General struct {
	I3StatusCMD string `json:"i3statuscmd"`
}

type Rule struct {
	// Unique rule ID
	ID string `json:"id"`

	// Name of the block, if set
	Name string `json:"name"`
	// Instance of the block, if set
	Instance string `json:"instance"`

	// Map buttons to commands which are executed when the event matches
	// key: X11 button ID (for example 1 to 3 for left/middle/right mouse button)
	// value: a command
	Actions map[string]string `json:"actions"`
}

func LoadConfig(fpath string) (*Config, error) {
	var cfg Config
	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
