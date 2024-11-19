package-darwin:
	# go install fyne.io/fyne/v2/cmd/fyne@latest
	fyne package --id ch.hacker.onboarder --os darwin --release --appVersion 1.0.0 --icon whisper-ui.png
	cp -v ffmpeg whisper-cpp ggml-medium.bin whisper-ui.app/Contents/Resources

clean:
	rm -rf whisper-ui whisper-ui.app
