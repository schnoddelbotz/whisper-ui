package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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
	mtxt := `Model information from [whisper-cpp github page](https://github.com/ggerganov/whisper.cpp/blob/master/models/README.md):
Models are multilingual unless the model name includes '.en'.
Models ending in '-q5_0' are [quantized](https://github.com/ggerganov/whisper.cpp/blob/master/README.md#quantization).
Models ending in '-tdrz' support local diarization (marking of speaker turns)
using [tinydiarize](https://github.com/akashmjn/tinydiarize).
More information about models is available
[upstream (openai/whisper)](https://github.com/openai/whisper#available-models-and-languages).`

	for _, m := range availableModels {
		mtxt += fmt.Sprintf("- %s (%s)\n", m.name, m.size)
	}
	modelText := widget.NewRichTextFromMarkdown(mtxt)
	modelText.Wrapping = fyne.TextWrapWord
	backToMainButton := widget.NewButton("Back to main menu", a.win_intro)
	backToMainButton.Enable()
	if len(a.installedModels) == 0 {
		backToMainButton.Disable()
	}

	top := container.NewVBox(
		widget.NewRichTextFromMarkdown("### Install a model for transcription"),
		widget.NewLabel("Model directory: "+getModelsDir()))
	bottom := container.NewVBox(
		a.install_model_chooser(),
		container.NewHBox(a.installModelButton, backToMainButton),
	)
	content := container.NewBorder(top, bottom, nil, nil, container.NewVScroll(modelText))
	a.window.SetContent(content)
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
	// https://github.com/ggerganov/whisper.cpp/blob/master/models/download-ggml-model.sh
	src := "https://huggingface.co/ggerganov/whisper.cpp"
	pfx := "resolve/main/ggml"
	if strings.Contains(a.installModel, "tdrz") {
		src = "https://huggingface.co/akashmjn/tinydiarize-whisper.cpp"
		pfx = "resolve/main/ggml"
	}
	url := fmt.Sprintf("%s/%s-%s.bin", src, pfx, a.installModel)
	tgt := filepath.Join(getModelsDir(), fmt.Sprintf("ggml-%s.bin", a.installModel))
	modelInfo := getModel(a.installModel)
	if modelInfo == nil {
		a.win_display_error(errors.New("invalid model name"), "Unexpected model name provided.")
		return
	}
	a.window.SetContent(
		container.NewVBox(
			widget.NewRichTextFromMarkdown(fmt.Sprintf(`### Downloading model "%s"
URL: %s

Size: %s

SHA: %s

Destination: %s`, a.installModel, url, modelInfo.size, modelInfo.sha, tgt)),
			layout.NewSpacer(),
			progressbar,
		),
	)
	if !ensureDirExists(getModelsDir()) {
		a.win_display_error(errors.New("models directory does not exist"), "... or is not writable")
		return
	}

	counter := &writeCounter{progress: boundProgress, sha: sha1.New()}
	size, err := httpHeadGetSize(url)
	if err != nil {
		a.win_display_error(err, "Failed to HTTP HEAD download size.")
	}
	counter.totalSize = size

	err = downloadFile(counter, tgt, url)
	if err != nil {
		a.win_display_error(err, "Download failed.")
	}

	if fmt.Sprintf("%x", counter.sha.Sum(nil)) != modelInfo.sha {
		a.win_display_error(errors.New("unexpected SHA sum for download"),
			fmt.Sprintf("%x does not match %s", counter.sha.Sum(nil), modelInfo.sha))
		return
	}

	a.installedModels = getModels()
	a.selectedModel = a.installModel
	go func() {
		a.window.SetContent(widget.NewLabel("Download and SHA verification successful."))
		time.Sleep(2 * time.Second)
		a.win_intro()
	}()

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
