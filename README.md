# whisper-ui

This is a quick PoC (read: hack) WIP for a [Whisper](https://github.com/ggerganov/whisper.cpp) UI.
Currently on macOS - but relying on [Fyne](https://fyne.io/) to be open for others, later, maybe.
The build / Makefile fetches a static ffmpeg build and builds whisper-cpp and puts it into
the .app bundle, including the 
[medium model](https://github.com/ggerganov/whisper.cpp/blob/master/models/README.md).

## building

```bash
git clone https://github.com/schnoddelbotz/whisper-ui.git
cd whisper-ui
make
```

## todo / issues

Status: Works for me (Sonoma/ARM, Monterey/x86_64), will possibly look into improvements.

- remove hard-coded language selection (german), provide chooser
- fix progress indicator / improve user feedback
- skip ffmpeg conversion if input is 16kHz WAV
- make it possible to use ffmpeg/whisper found in PATH, "unbundled" build
- prevent multiple clicks on "select input file"
- Sequioa `Bad CPU Type in Executable`: `softwareupdate --install-rosetta`?

## license

MIT