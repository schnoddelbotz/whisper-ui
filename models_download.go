package main

import (
	"hash"
	"io"
	"net/http"
	"os"
	"strconv"

	"fyne.io/fyne/v2/data/binding"
	"golang.org/x/sys/unix"
)

type writeCounter struct {
	written   uint64
	totalSize int
	progress  binding.ExternalFloat
	sha       hash.Hash
}

// based on https://www.golangcode.com/download-a-file-with-progress/

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

func httpHeadGetSize(url string) (int, error) {
	headResp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer headResp.Body.Close()
	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))
	if err != nil {
		return 0, err
	}
	return size, nil
}

func ensureDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false
		}
	}
	return unix.Access(path, unix.W_OK) == nil
}
