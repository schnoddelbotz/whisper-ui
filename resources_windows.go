package main

import (
	"os"
	"path/filepath"
)

func getResources() *resources {
	// Whisper and ffmpeg expected in same direcotry as whisper-ui.exe
	dir := getResourcesDir()
	rsrc := &resources{
		ffmpeg:     dir + `\ffmpeg.exe`,
		whispercpp: dir + `\whisper-cli.exe`,
		tmpfile:    getTempFileName(),
	}
	return rsrc
}

func getResourcesDir() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}

func getModelsDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, "whisper-ui")
}
