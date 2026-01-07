package main

import (
	"os"
	"path/filepath"
)

func getResources() *resources {
	// Whisper, its model and ffmpeg reside in .app bundle on macOS.
	dir := getResourcesDir()
	rsrc := &resources{
		ffmpeg:     dir + "/ffmpeg",
		whispercpp: dir + "/whisper-cli",
		tmpfile:    getTempFileName(),
	}
	return rsrc
}

func getResourcesDir() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex) + "/../Resources"
}

func getModelsDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, "Library/Application Support/whisper-ui")
}
