package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
)

// WrapI3Status starts i3status in a separate process,
// enables I3 click events and then prcoeeds to forward i3status stdout to our stdout.
func WrapI3Status(i3statuscmd string) error {
	os.Stdout.Write([]byte("{\"version\":1,\"click_events\":true}\n"))

	cmd := exec.Command(i3statuscmd)
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = os.Stderr
	go func() {
		r := bufio.NewReader(pr)
		//skip header line which is missing the click_events flag
		_, _ = r.ReadBytes('\n')
		for {
			// redirect stdout
			line, err := r.ReadBytes('\n')
			if err != nil {
				log.Println("[ERROR]:", err)
				continue
			}
			os.Stdout.Write(line)
		}
	}()
	log.Println("[INFO] Starting", i3statuscmd)
	return cmd.Run()
}
