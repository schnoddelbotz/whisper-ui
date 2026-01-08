# whisper-ui

whisper-ui is a very simple [whisper-cpp](https://github.com/ggerganov/whisper.cpp) GUI wrapper.
It enables local/offline usage of [Whisper](https://openai.com/index/whisper/) to transcribe audio
or video input files to plain text. whisper-ui uses [Fyne](https://fyne.io/) to build the UI.

[Download a whisper-ui release](https://github.com/schnoddelbotz/whisper-ui/releases),
available for macOS, Windows and Ubuntu. Releases bundle a whisper-cpp and
[ffmpeg](https://www.ffmpeg.org/download.html) executable to free users from
any further setup/compilation tasks.

[Models for whisper-cpp](https://github.com/ggerganov/whisper.cpp/blob/master/models/README.md)
can be downloaded using the whisper-ui app.

## Notes - macOS

Note that the releases built via github [workflow](.github/workflows/release.yaml) are not signed.
For macOS, this means you have to remove quarantine flag (using `xattr -dr com.apple.quarantine whisper-ui.app`).

To build whisper-ui from source:

```bash
git clone https://github.com/schnoddelbotz/whisper-ui.git
cd whisper-ui
make build-darwin
```

The build / Makefile fetches a [static ffmpeg build](https://evermeet.cx/ffmpeg/),
builds whisper-cpp and puts both into the .app bundle.

## Notes - Linux

The Linux release contains a Makefile. It can be used to install or uninstall
whisper-ui as Desktop app. For a system-wide install, use `sudo make install`.
For an installation for the current user, use `make user-install` -- note that
for this to work, your `PATH` must contain `~/.local/bin`.

The Linux build expects ffmpeg to be installed via package manger / in `$PATH`.
If not already present, use `sudo apt install ffmpeg`.

To build whisper-ui from source:

```bash
git clone https://github.com/schnoddelbotz/whisper-ui.git
cd whisper-ui
make build zip-linux
```

## Notes - Windows

Like the macOS release, whisper-ui.exe is not signed - accordingly
Defender will report a "unrecognized app" and clicking "More info"
and then "Run anyway" will be required when starting it for the first time.

For building whisper-ui from source, see the
[release.yaml](.github/workflows/release.yaml)'s windows section.

## todo / issues

Status: Toy project, works for me (tested on Tahoe/ARM, Windows 11 & Ubuntu 24.04 on x86_64).

- maybe add [whisper-stream](https://github.com/ggml-org/whisper.cpp?tab=readme-ov-file#real-time-audio-input-example)?
- maybe embed [FyneTerm](https://github.com/fyne-io/terminal) to show
  [confidence color-coding](https://github.com/ggml-org/whisper.cpp?tab=readme-ov-file#confidence-color-coding)?
  This would also solve the next issue:
- more verbose output/progress feedback from ffmpeg and whisper-cpp?
- could skip ffmpeg conversion if input is already a 16kHz WAV
- Sequioa/ARM using x86 build gives `Bad CPU Type in Executable`? `softwareupdate --install-rosetta`.
- option to "open output upon completion" may not work on Windows, same for Markdown link to file
- it would be nice to be able to [select multiple files](https://github.com/fyne-io/fyne/issues/1082) for conversion

## license

MIT
