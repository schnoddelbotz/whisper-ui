package main

import (
	"log"
	"os"
)

type resources struct {
	ffmpeg     string
	whispercpp string
	model      string
	tmpfile    string
}

func getTempFileName() string {
	// just interested in a name in $TMP that does not exist.
	// CreateTemp ensures this ... is there a way to avoid creation to just get name?
	f, err := os.CreateTemp("", "whisper-ui")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return f.Name() + ".wav"
}
