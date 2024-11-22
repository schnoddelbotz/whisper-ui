# whisper-ui

This is a quick PoC (read: hack) WIP for a [whisper-cpp](https://github.com/ggerganov/whisper.cpp) UI.
It enables local/offline usage of [Whisper](https://openai.com/index/whisper/) to transcribe audio
or video input files to plain text.
Currently on macOS - but relying on [Fyne](https://fyne.io/) to be open for others, later, maybe.

## building - macOS

```bash
git clone https://github.com/schnoddelbotz/whisper-ui.git
cd whisper-ui
make
```

The build / Makefile fetches a static ffmpeg build, builds whisper-cpp and puts both into
the .app bundle. Models can be downloaded using the app.

## todo / issues

Status: Works for me (Sonoma/ARM, Monterey/x86_64), will possibly look into improvements.

- build/release as github action (go-releaser?)
- skip ffmpeg conversion if input is 16kHz WAV
- make it possible to use ffmpeg/whisper found in PATH, "unbundled" build
- Sequioa/ARM using x86 build gives `Bad CPU Type in Executable`? `softwareupdate --install-rosetta`.
- ahem, tests?

## license

MIT
