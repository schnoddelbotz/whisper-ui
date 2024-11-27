# whisper-ui

This is a quick PoC (read: hack) for a [whisper-cpp](https://github.com/ggerganov/whisper.cpp) UI.
It enables local/offline usage of [Whisper](https://openai.com/index/whisper/) to transcribe audio
or video input files to plain text. whisper-ui uses [Fyne](https://fyne.io/) to build the UI.

[Download a whisper-ui release](https://github.com/schnoddelbotz/whisper-ui/releases),
available for macOS and Windows and Ubuntu.

Note that the releases built via github [workflow](.github/workflows/release.yaml) are not signed.
For macOS, this means you have to remove quarantine flag (using `xattr -d com.apple.quarantine ...`).

## building - macOS

```bash
git clone https://github.com/schnoddelbotz/whisper-ui.git
cd whisper-ui
make build-darwin
```

The build / Makefile fetches a static ffmpeg build, builds whisper-cpp and puts both into
the .app bundle. Models can be downloaded using the app.

## building - Linux

```bash
git clone https://github.com/schnoddelbotz/whisper-ui.git
cd whisper-ui
make build zip-linux
```

The Linux build expects ffmpeg to be installed via package manger / in `$PATH`.
The .zip produced contains a Makefile - run `make user-install`.

## building - Windows

See the [release.yaml](.github/workflows/release.yaml)'s windows section.

## todo / issues

Status: Works for me (Sonoma/ARM, Monterey/x86_64), will possibly look into improvements.

- skip ffmpeg conversion if input is already a 16kHz WAV
- make it possible to use ffmpeg/whisper found in PATH, "unbundled" build
- more verbose output/feedback from ffmpeg and whisper-cpp?
- Sequioa/ARM using x86 build gives `Bad CPU Type in Executable`? `softwareupdate --install-rosetta`.
- ahem, tests?

## license

MIT
