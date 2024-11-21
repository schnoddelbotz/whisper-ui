package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type application struct {
	app                fyne.App
	window             fyne.Window
	openWhenDone       binding.Bool
	selectFileButton   *widget.Button
	selectLanguage     *widget.Select
	selectedLanguage   string
	installedModels    []string
	selectModel        *widget.Select
	selectedModel      string
	installModel       string
	installModelButton *widget.Button
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
	a.selectFileButton = widget.NewButton("Select input file", a.filechooser_input)
	a.selectedLanguage = "de"
	a.selectLanguage = widget.NewSelect([]string{"de", "en"}, func(value string) {
		a.selectedLanguage = value
	})

	a.installedModels = getModels()
	if len(a.installedModels) > 0 {
		a.selectedModel = a.installedModels[0]
	}

	a.win_intro()
	a.window.ShowAndRun()
}

func (a *application) win_intro() {
	if len(a.installedModels) == 0 {
		a.win_install_model()
		return
	}
	a.selectModel = widget.NewSelect(a.installedModels, func(value string) {
		a.selectedModel = value
	})
	a.selectLanguage.SetSelected(a.selectedLanguage)
	a.selectFileButton.Enable()
	a.selectModel.SetSelected(a.selectedModel)
	a.window.SetContent(
		container.NewVBox(
			widget.NewLabel("This tool converts audio/video to text (offline)."),
			widget.NewCheckWithData("Open file upon completion", a.openWhenDone),
			container.NewHBox(widget.NewLabel("Language"), a.selectLanguage),
			container.NewHBox(widget.NewLabel("Model"), a.selectModel,
				widget.NewButton("Add model", a.win_install_model)),
			layout.NewSpacer(),
			a.selectFileButton,
		),
	)
	a.window.Canvas().Focus(a.selectFileButton)
}

func (a *application) win_install_model() {
	a.installModelButton = widget.NewButton("Install selected model", a.win_downloading_model)
	a.installModelButton.Disable()
	mtxt := ""
	// mtxt := "Models are multilingual unless the model name includes `.en`."
	// mtxt += "Models ending in `-q5_0` are [quantized](../README.md#quantization)."
	// mtxt += "Models ending in `-tdrz` support local diarization (marking of speaker turns) "
	// mtxt += "using [tinydiarize](https://github.com/akashmjn/tinydiarize). "
	// mtxt += "More information about models is available "
	// mtxt += "[upstream (openai/whisper)](https://github.com/openai/whisper#available-models-and-languages).\n"
	// ^ adding this temporarily stretches window height, making bottom button inaccessible. why?
	for _, m := range availableModels {
		mtxt += fmt.Sprintf("- %s (%s)\n", m.name, m.size)
	}
	modelText := widget.NewRichTextFromMarkdown(mtxt)
	modelText.Wrapping = fyne.TextWrapWord
	backToMainButton := widget.NewButton("Back to main menu", a.win_intro)
	a.window.SetContent(
		container.NewVBox(
			widget.NewRichTextFromMarkdown("### Install a model for transcription"),
			widget.NewLabel("Model directory: "+getModelsDir()),
			modelText,
			a.install_model_chooser(),
			layout.NewSpacer(),
			container.NewHBox(a.installModelButton, backToMainButton),
		),
	)
	backToMainButton.Enable()
	if len(a.installedModels) == 0 {
		backToMainButton.Disable()
	}
}

func (a *application) install_model_chooser() *widget.Select {
	installable := []string{}
	for _, m := range availableModels {
		if modelIsInstalled(m.name, a.installedModels) {
			continue
		}
		installable = append(installable, m.name)
	}
	return widget.NewSelect(installable, func(value string) {
		a.installModel = value
		a.installModelButton.Enable()
	})
}

func (a *application) win_downloading_model() {
	var progress float64
	boundProgress := binding.BindFloat(&progress)
	progressbar := widget.NewProgressBarWithData(boundProgress)
	src := "https://huggingface.co/ggerganov/whisper.cpp"
	pfx := "resolve/main/ggml"
	// todo: tinything case
	// https://github.com/ggerganov/whisper.cpp/blob/master/models/download-ggml-model.sh#L87
	url := fmt.Sprintf("%s/%s-%s.bin", src, pfx, a.installModel)
	tgt := filepath.Join(getModelsDir(), fmt.Sprintf("ggml-%s.bin", a.installModel))
	a.window.SetContent(
		container.NewVBox(
			widget.NewRichTextFromMarkdown(fmt.Sprintf(`### Downloading model %s
URL: %s

Size: xxx,

SHA: xxx,

Destination: %s`, a.installModel, url, tgt)),
			layout.NewSpacer(),
			progressbar,
		),
	)
	if !ensureDirExists(getModelsDir()) {
		a.win_display_error(errors.New("models directory does not exist"), "... or is not writable")
		return
	}

	counter := &writeCounter{progress: boundProgress}
	size, err := httpHeadGetSize(url)
	if err != nil {
		a.win_display_error(err, "Failed to HTTP HEAD download size.")
	}
	counter.totalSize = size

	err = downloadFile(counter, tgt, url)
	if err != nil {
		a.win_display_error(err, "Download failed.")
	}

	a.installedModels = getModels()
	a.selectedModel = a.installModel
	a.win_intro()
}

func (a *application) filechooser_input() {
	a.selectFileButton.Disable()
	dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		if f == nil {
			a.selectFileButton.Enable()
			return
		}
		a.win_convert(f.URI().Path())
	}, a.window)
}

func (a *application) win_convert(f string) {
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
		a.win_display_error(err, string(cmdout))
		return
	}

	statusText.SetText("Transcribing using Whisper ...")
	model := filepath.Join(getModelsDir(), fmt.Sprintf("ggml-%s.bin", a.selectedModel))
	cmd = exec.Command(rsrc.whispercpp, "-l", a.selectedLanguage, "-m", model, "-f", rsrc.tmpfile, "-otxt", "-of", f) // auto-appends .txt
	cmdout, err = cmd.CombinedOutput()
	os.Remove(rsrc.tmpfile)
	if err != nil {
		a.win_display_error(err, string(cmdout))
		return
	}

	a.win_success(f + ".txt")
}

func (a *application) win_success(filename string) {
	a.window.SetContent(
		container.NewVBox(
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

func (a *application) win_display_error(err error, msg string) {
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
