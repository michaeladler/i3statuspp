package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("share/config.json")
	assert.NoError(t, err)
	assert.Equal(t, "i3status", cfg.General.I3StatusCMD)
	assert.Len(t, cfg.Rules, 1)
	assert.Equal(t, "1", cfg.Rules[0].ID)
	assert.Equal(t, "volume", cfg.Rules[0].Name)
	assert.Equal(t, "pulse:alsa_output.pci-0000_00_1b.0.analog-stereo.Master.0", cfg.Rules[0].Instance)
	assert.Equal(t, "pavucontrol", cfg.Rules[0].Actions["1"])
	assert.Equal(t, "pamixer --sink alsa_output.pci-0000_00_1b.0.analog-stereo --allow-boost --increase 1", cfg.Rules[0].Actions["4"])
	assert.Equal(t, "pamixer --sink alsa_output.pci-0000_00_1b.0.analog-stereo --allow-boost --decrease 1", cfg.Rules[0].Actions["5"])
}
