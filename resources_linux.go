package main

import (
	"log"
	"os"
)

func getResources() *resources {
	dir := getResourcesDir()
	rsrc := &resources{
		ffmpeg:     "ffmpeg",                 // expect in PATH for now
		whispercpp: "whisper-cpp",            // expect in PATH for now
		model:      dir + "/ggml-medium.bin", // ~/.whisper-ui/ggml-medium.bin
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
