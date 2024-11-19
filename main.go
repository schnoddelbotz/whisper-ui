package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type application struct {
	app    fyne.App
	window fyne.Window
}

func main() {
	a := application{
		app: app.New(),
	}

	a.window = a.app.NewWindow("whisper-ui")

	a.window.Resize(fyne.NewSize(600, 400))
	a.window.SetFixedSize(true)

	a.win_intro()
	a.window.ShowAndRun()
}

func (a *application) win_intro() {
	next_button := widget.NewButton("Select input file", a.select_input)
	a.window.SetContent(
		container.NewVBox(
			widget.NewLabel("This tool converts audio/video to text."),
			widget.NewLabel("Select an input file for conversion."),
			layout.NewSpacer(),
			next_button,
		),
	)
}

func (a *application) select_input() {
	dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
		saveFile := "NoFileYet"
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		if f == nil {
			return
		}
		saveFile = f.URI().Path()
		a.convert(saveFile)
	}, a.window)
}

func (a *application) convert(f string) {
	infinite := widget.NewProgressBarInfinite()
	// handle invalid file type ...
	a.window.SetContent(
		container.NewVBox(
			widget.NewLabel("Converting "+f),
			widget.NewLabel("Please wait ..."),
			layout.NewSpacer(),
			infinite, // why is this not animated...?
		),
	)
	// start conversion
}
