package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type application struct {
	app                fyne.App
	window             fyne.Window
	openWhenDone       binding.Bool
	translate          binding.Bool
	outputFormat       string
	selectFileButton   *widget.Button
	selectedLanguage   string
	installedModels    []string
	selectedModel      string
	installModel       string
	installModelButton *widget.Button
}

var languages = []string{
	// https://github.com/ggerganov/whisper.cpp/blob/master/src/whisper.cpp#L246
	"auto", "en", "de", "es", "ru", "ko", "fr", "ja", "pt", "tr", "pl", "ca", "nl",
	"ar", "sv", "it", "id", "hi", "fi", "vi", "he", "uk", "el", "ms", "cs", "ro",
	"da", "hu", "ta", "no", "th", "ur", "hr", "bg", "lt", "la", "mi", "ml", "cy",
	"sk", "te", "fa", "lv", "bn", "sr", "az", "sl", "kn", "et", "mk", "br", "eu",
	"is", "hy", "ne", "mn", "bs", "kk", "sq", "sw", "gl", "mr", "pa", "si", "km",
	"sn", "yo", "so", "af", "oc", "ka", "be", "tg", "sd", "gu", "am", "yi", "lo",
	"uz", "fo", "ht", "ps", "tk", "nn", "mt", "sa", "lb", "my", "bo", "tl", "mg",
	"as", "tt", "haw", "ln", "ha", "ba", "jw", "su", "yue", "zh",
}

func main() {
	a := application{
		app: app.New(),
	}

	a.window = a.app.NewWindow("whisper-ui")

	a.window.Resize(fyne.NewSize(700, 500))
	a.window.SetFixedSize(true)

	a.openWhenDone = binding.NewBool()
	a.translate = binding.NewBool()
	a.selectFileButton = widget.NewButton("Select input file", a.windowInputFileChooser)
	a.selectedLanguage = "de"
	a.outputFormat = "txt"

	a.installedModels = getModels()
	if len(a.installedModels) > 0 {
		a.selectedModel = a.installedModels[0]
	}

	a.windowMain()
	a.window.ShowAndRun()
}

func (a *application) windowMain() {
	if len(a.installedModels) == 0 {
		a.windowInstallModel()
		return
	}

	// https://github.com/fyne-io/fyne/issues/2836 - No binding support for select :(
	selectLanguage := widget.NewSelect(languages, func(value string) {
		a.selectedLanguage = value
	})
	selectLanguage.SetSelected(a.selectedLanguage)

	selectModel := widget.NewSelect(a.installedModels, func(value string) {
		a.selectedModel = value
	})
	selectModel.SetSelected(a.selectedModel)

	selectFormat := widget.NewSelect([]string{"txt", "vtt", "srt", "lrc", "wts"}, func(value string) {
		a.outputFormat = value
	})
	selectFormat.SetSelected(a.outputFormat)

	a.selectFileButton.Enable()
	about := "<br><br>\n### About\nwhisper-ui solely provides a very limited user interface for\n"
	about += "- [whsiper-cpp](https://github.com/ggerganov/whisper.cpp), a"
	about += " [Whisper](https://github.com/openai/whisper) C++ binary implementation, used for transcription\n"
	about += "- [ffmpeg](https://www.ffmpeg.org/), the universal audio/video converter\n\n"
	about += "<br>To ease usage, whisper-ui releases bundle both tools as executable binary.\n\n"
	about += "### Issues\n[Report a whisper-ui issue](https://github.com/schnoddelbotz/whisper-ui/issues).\n"
	a.window.SetContent(
		container.NewVBox(
			widget.NewRichTextFromMarkdown("### Use whisper-ui to transcribe audio/video to text (offline)."),
			widget.NewCheckWithData("Open output file upon completion", a.openWhenDone),
			widget.NewCheckWithData("Translate from source language to english", a.translate),
			container.NewHBox(widget.NewLabel("Output file format"), selectFormat),
			container.NewHBox(widget.NewLabel("Spoken language"), selectLanguage),
			container.NewHBox(widget.NewLabel("Model"), selectModel,
				widget.NewButton("Add model", a.windowInstallModel)),
			widget.NewRichTextFromMarkdown(about),
			layout.NewSpacer(),
			container.NewHBox(
				a.selectFileButton,
				layout.NewSpacer(),
				widget.NewButton("Quit", a.app.Quit),
			),
		),
	)
	a.window.Canvas().Focus(a.selectFileButton)
}

func (a *application) windowInstallModel() {
	a.installModelButton = widget.NewButton("Install selected model", a.windowDownloadingModel)
	a.installModelButton.Disable()
	mtxt := `Model information from [whisper-cpp github page](https://github.com/ggerganov/whisper.cpp/blob/master/models/README.md):
Models are multilingual unless the model name includes '.en'.
Models ending in '-q5_0' are [quantized](https://github.com/ggerganov/whisper.cpp/blob/master/README.md#quantization).
Models ending in '-tdrz' support local diarization (marking of speaker turns)
using [tinydiarize](https://github.com/akashmjn/tinydiarize).
More information about models is available
[upstream (openai/whisper)](https://github.com/openai/whisper#available-models-and-languages).

`

	for _, m := range availableModels {
		mtxt += fmt.Sprintf("- %s (%s)\n", m.name, m.size)
	}
	modelText := widget.NewRichTextFromMarkdown(mtxt)
	modelText.Wrapping = fyne.TextWrapWord
	backToMainButton := widget.NewButton("Back to main menu", a.windowMain)
	backToMainButton.Enable()
	if len(a.installedModels) == 0 {
		backToMainButton.Disable()
	}

	top := container.NewVBox(
		widget.NewRichTextFromMarkdown("### Install a model for transcription"),
		widget.NewLabel("Model directory: "+getModelsDir()))
	bottom := container.NewVBox(
		a.installModelSelect(),
		container.NewHBox(
			a.installModelButton,
			layout.NewSpacer(),
			backToMainButton,
			widget.NewButton("Quit", a.app.Quit)),
	)
	content := container.NewBorder(top, bottom, nil, nil, container.NewVScroll(modelText))
	a.window.SetContent(content)
}

func (a *application) installModelSelect() *widget.Select {
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

func (a *application) windowDownloadingModel() {
	counter := NewWriteCounter()
	url := getURLForModel(a.installModel)
	tgt := filepath.Join(getModelsDir(), fmt.Sprintf("ggml-%s.bin", a.installModel))
	modelInfo := getModel(a.installModel)
	if modelInfo == nil {
		a.windowDisplayError(errors.New("invalid model name"), "Unexpected model name provided.")
		return
	}
	infoMarkdown := "### Downloading model '%s'\nURL: %s\n\nSize: %s\n\nSHA: %s\n\nDestination: %s"
	a.window.SetContent(
		container.NewVBox(
			widget.NewRichTextFromMarkdown(fmt.Sprintf(infoMarkdown, a.installModel, url, modelInfo.size, modelInfo.sha, tgt)),
			layout.NewSpacer(),
			widget.NewProgressBarWithData(counter.progress),
		),
	)
	if !ensureDirExists(getModelsDir()) {
		a.windowDisplayError(errors.New("models directory does not exist"), "... or is not writable")
		return
	}

	size, err := httpHeadGetSize(url)
	if err != nil {
		a.windowDisplayError(err, "Failed to HTTP HEAD download size.")
	}
	counter.totalSize = size

	err = downloadFile(counter, tgt, url)
	if err != nil {
		a.windowDisplayError(err, "Download failed.")
	}

	if fmt.Sprintf("%x", counter.sha.Sum(nil)) != modelInfo.sha {
		a.windowDisplayError(errors.New("unexpected SHA sum for download"),
			fmt.Sprintf("%x does not match %s", counter.sha.Sum(nil), modelInfo.sha))
		return
	}

	a.installedModels = getModels()
	a.selectedModel = a.installModel
	go func() {
		a.window.SetContent(widget.NewLabel("Download and SHA verification successful."))
		time.Sleep(2 * time.Second)
		a.windowMain()
	}()
}

func (a *application) windowInputFileChooser() {
	a.selectFileButton.Disable()
	d := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		if f == nil {
			a.selectFileButton.Enable()
			return
		}
		a.windowConverting(f.URI().Path())
	}, a.window)
	d.SetFilter(storage.NewExtensionFileFilter([]string{".mp3", ".mp4", ".mov", ".mpg", ".wav", ".m4a", ".ogg", ".avi", ".mkv"}))
	d.Show()
}

func (a *application) windowConverting(f string) {
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
		a.windowDisplayError(err, string(cmdout))
		return
	}

	statusText.SetText("Transcribing using Whisper ...")
	translate := ""
	if doit, _ := a.translate.Get(); doit {
		translate = "-tr"
	}
	model := filepath.Join(getModelsDir(), fmt.Sprintf("ggml-%s.bin", a.selectedModel))
	cmd = exec.Command(rsrc.whispercpp, "-l", a.selectedLanguage, "-m", model, "-f", rsrc.tmpfile, "-o"+a.outputFormat, "-of", f, translate) // auto-appends .txt
	cmdout, err = cmd.CombinedOutput()
	os.Remove(rsrc.tmpfile)
	if err != nil {
		a.windowDisplayError(err, string(cmdout))
		return
	}

	a.windowConversionSuccess(f + "." + a.outputFormat)
}

func (a *application) windowConversionSuccess(filename string) {
	transcribeNextButton := widget.NewButton("Transcribe another file", a.windowMain)
	fileURL := url.URL{Path: filename}
	fileLink := fmt.Sprintf("[%s](file://%s)", filename, fileURL.String())
	a.window.SetContent(
		container.NewVBox(
			widget.NewRichTextFromMarkdown("Done: "+fileLink),
			layout.NewSpacer(),
			container.NewHBox(
				transcribeNextButton,
				layout.NewSpacer(),
				widget.NewButton("Quit", a.app.Quit),
			),
		),
	)
	a.window.Canvas().Focus(transcribeNextButton)
	if doit, _ := a.openWhenDone.Get(); doit {
		// todo: Linux xdg-open, Windows rundll32
		// see also https://github.com/golang/go/issues/32456
		exec.Command("open", filename).Start()
	}
}

func (a *application) windowDisplayError(err error, msg string) {
	quit_button := widget.NewButton("Quit", a.app.Quit)

	top := widget.NewLabel("Error: " + err.Error())
	bottom := container.NewHBox(
		widget.NewButton("Back to main menu", a.windowMain),
		layout.NewSpacer(),
		widget.NewButton("Quit", a.app.Quit),
	)
	content := container.NewBorder(top, bottom, nil, nil, container.NewVScroll(widget.NewLabel(msg)))
	a.window.SetContent(content)

	a.window.Canvas().Focus(quit_button)
}
