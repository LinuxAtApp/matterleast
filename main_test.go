package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	stdout := os.Stdout
	// Credit to Evan Shaw for this way to capture stdout without creating
	// a new file. See:
	// http://stackoverflow.com/a/10476304/7106084
	readEnd, writeEnd, err := os.Pipe()
	if err != nil {
		t.Skip("Unable to create pipe to replace stdout")
	}
	os.Stdout = writeEnd
	// Send anything written on stdout back on the contents channel as a string
	contents := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, readEnd)
		contents <- buf.String()
	}()
	main() //print to stdout
	os.Stdout = stdout
	writeEnd.Close() // tell our goroutine that we're done writing

	stringContents := <-contents // wait for goroutine to give us output
	if !strings.Contains(stringContents, "matterleast") {
		t.Error("Stdout incorrect; expected the word 'matterleast'")
	}
}
