package main

import (
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type application struct {
	app              fyne.App
	window           fyne.Window
	openWhenDone     binding.Bool
	selectFileButton *widget.Button
	selectLanguage   *widget.Select
	selectedLanguage string
}

func main() {
	a := application{
		app: app.New(),
	}

	a.window = a.app.NewWindow("whisper-ui")

	a.window.Resize(fyne.NewSize(700, 500))
	a.window.SetFixedSize(true)

	a.openWhenDone = binding.NewBool()
	a.openWhenDone.Set(true)
	a.selectFileButton = widget.NewButton("Select input file", a.select_input)
	a.selectedLanguage = "de"
	a.selectLanguage = widget.NewSelect([]string{"de", "en"}, func(value string) {
		a.selectedLanguage = value
	})

	a.win_intro()
	a.window.ShowAndRun()
}

func (a *application) win_intro() {
	a.selectLanguage.SetSelected(a.selectedLanguage)
	a.selectFileButton.Enable()
	a.window.SetContent(
		container.NewVBox(
			widget.NewLabel("This tool converts audio/video to text (offline)."),
			widget.NewCheckWithData("Open file upon completion", a.openWhenDone),
			container.NewHBox(widget.NewLabel("Language"), a.selectLanguage),
			layout.NewSpacer(),
			a.selectFileButton,
		),
	)
	a.window.Canvas().Focus(a.selectFileButton)
}

func (a *application) select_input() {
	a.selectFileButton.Disable()
	dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
		saveFile := "NoFileYet"
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		if f == nil {
			a.selectFileButton.Enable()
			return
		}
		saveFile = f.URI().Path()
		a.convert(saveFile)
	}, a.window)
}

func (a *application) convert(f string) {
	// todo: select language, improve progress feedback, ... CLEAN UP (eg. os.TempDir, stderr...).
	// handle invalid file type ...
	progressBar := widget.NewProgressBarInfinite()
	statusText := widget.NewLabel("Converting input file to 16 kHz .WAV")
	a.window.SetContent(
		container.NewVBox(
			widget.NewLabel("Converting "+f),
			widget.NewLabel("Please wait."),
			statusText,
			layout.NewSpacer(),
			progressBar,
		),
	)
	progressBar.Refresh()

	// start conversion
	rsrc := getResources()
	cmd := exec.Command(rsrc.ffmpeg, "-i", f, "-acodec", "pcm_s16le", "-ac", "1", "-ar", "16000", rsrc.tmpfile)
	cmdout, err := cmd.CombinedOutput()
	if err != nil {
		a.display_error(err, string(cmdout))
		return
	}

	statusText.SetText("Transcribing using Whisper ...")
	cmd = exec.Command(rsrc.whispercpp, "-l", a.selectedLanguage, "-m", rsrc.model, "-f", rsrc.tmpfile, "-otxt", "-of", f) // auto-appends .txt
	cmdout, err = cmd.CombinedOutput()
	os.Remove(rsrc.tmpfile)
	if err != nil {
		a.display_error(err, string(cmdout))
		return
	}

	a.display_success(f + ".txt")
}

func (a *application) display_success(filename string) {
	a.window.SetContent(
		container.NewVBox(
			//widget.NewLabel("Done: "+filename), // todo: link me (markdown...?)
			widget.NewRichTextFromMarkdown("Done: ["+filename+"](file://"+filename+")"),
			layout.NewSpacer(),
			container.NewHBox(
				widget.NewButton("Transcribe another file", a.win_intro),
				widget.NewButton("Quit", a.app.Quit),
			),
		),
	)
	if doit, _ := a.openWhenDone.Get(); doit {
		// todo: Linux xdg-open, Windows rundll32
		// see also https://github.com/golang/go/issues/32456
		exec.Command("open", filename).Start()
	}
}

func (a *application) display_error(err error, msg string) {
	quit_button := widget.NewButton("Quit", a.app.Quit)
	a.window.SetContent(
		container.NewVBox(
			widget.NewLabel("Error: "+err.Error()),
			widget.NewLabel(msg),
			layout.NewSpacer(),
			quit_button,
		),
	)
	a.window.Canvas().Focus(quit_button)
}
