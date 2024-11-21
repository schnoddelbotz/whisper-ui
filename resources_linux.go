package main

import (
	"log"
	"os"
	"path/filepath"
)

// TODO / WIP / TBD ... Incomplete.

func getResources() *resources {
	rsrc := &resources{
		ffmpeg:     "ffmpeg",      // expect in PATH for now
		whispercpp: "whisper-cpp", // expect in PATH for now
		tmpfile:    getTempFileName(),
	}
	return rsrc
}

func getResourcesDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err) // not useful as GUI app ... improve.
	}
	return dirname + "/.whisper-ui"
}

func getModelsDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".whisper-ui")
}
