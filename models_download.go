package main

import (
	"crypto/sha1"
	"hash"
	"io"
	"net/http"
	"os"
	"strconv"

	"fyne.io/fyne/v2/data/binding"
)

type writeCounter struct {
	written   uint64
	totalSize int
	progress  binding.ExternalFloat
	sha       hash.Hash
}

// based on https://www.golangcode.com/download-a-file-with-progress/
// maybe switch to https://github.com/ggerganov/whisper.cpp/blob/master/bindings/go/examples/go-model-download/main.go

func NewWriteCounter() *writeCounter {
	var progress float64
	boundProgress := binding.BindFloat(&progress)
	return &writeCounter{progress: boundProgress, sha: sha1.New()}
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.sha.Write(p)
	wc.written += uint64(n)
	wc.progress.Set(float64(wc.written) / float64(wc.totalSize))
	return n, nil
}

func downloadFile(counter *writeCounter, filepath string, url string) error {
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	counter.totalSize, _ = strconv.Atoi(resp.Header.Get("Content-Length"))

	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

func ensureDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false
		}
	}
	return true // FIXME -windows unix.Access(path, unix.W_OK) == nil
}
