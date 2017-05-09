package main

import (
	"bytes"
	"io"
	"os"
	"path"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	stdout := os.Stdout
	tmpStdout, err := os.Create(path.Join(os.TempDir(), "test-stdout"))
	if err != nil {
		t.Skip("Unable to create temporary test file")
	}
	os.Stdout = tmpStdout
	main()
	os.Stdout = stdout
	contents := make([]byte, 256)
	_, err = io.ReadFull(tmpStdout, contents)
	if err != nil {
		t.Skip("Unable to read output buffer")
	}
	stringContents := bytes.NewBuffer(contents).String()
	if !strings.Contains(stringContents, "matterleast") {
		t.Error("Stdout incorrect; expected the word 'matterleast'")
	}
}
