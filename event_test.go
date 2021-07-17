package main

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadEvents(t *testing.T) {

	inputs := []string{
		`{"name":"tztime","instance":"local","button":1,"modifiers":[],"x":1839,"y":8,"relative_x":31,"relative_y":8,"output_x":1839,"output_y":8,"width":112,"height":20}`,
		`,{"name":"tztime","instance":"local","button":1,"modifiers":[],"x":1839,"y":8,"relative_x":31,"relative_y":8,"output_x":1839,"output_y":8,"width":112,"height":20}`,
		`[{"name":"tztime","instance":"local","button":1,"modifiers":[],"x":1839,"y":8,"relative_x":31,"relative_y":8,"output_x":1839,"output_y":8,"width":112,"height":20}`,
	}
	for _, s := range inputs {

		eventChan := make(chan ClickEvent, 1)
		rd := bytes.NewReader([]byte(s))
		go func() {
			_ = ReadEvents(bufio.NewReader(rd), eventChan)
		}()

		ev := <-eventChan
		assert.Equal(t, "tztime", ev.Name)
		assert.Equal(t, "local", ev.Instance)
		assert.Equal(t, 1, ev.Button)
	}

}

func TestReadEventsBatch(t *testing.T) {

	input := `{"name":"first","instance":"local","button":1,"modifiers":[],"x":1839,"y":8,"relative_x":31,"relative_y":8,"output_x":1839,"output_y":8,"width":112,"height":20}
, {"name":"second","instance":"local","button":1,"modifiers":[],"x":1839,"y":8,"relative_x":31,"relative_y":8,"output_x":1839,"output_y":8,"width":112,"height":20}`

	eventChan := make(chan ClickEvent, 2)
	rd := bytes.NewReader([]byte(input))
	go func() {
		_ = ReadEvents(bufio.NewReader(rd), eventChan)
	}()

	ev := <-eventChan
	assert.Equal(t, "first", ev.Name)
	ev = <-eventChan
	assert.Equal(t, "second", ev.Name)
}
