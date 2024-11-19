package-darwin: clean ffmpeg whisper.cpp
	fyne version || go install fyne.io/fyne/v2/cmd/fyne@latest
	fyne package --id ch.hacker.onboarder --os darwin --release --appVersion 1.0.0 --icon whisper-ui.png
	ln whisper.cpp/main whisper-ui.app/Contents/Resources/whisper-cpp
	ln whisper.cpp/models/ggml-medium.bin whisper-ui.app/Contents/Resources/ggml-medium.bin
	ln ffmpeg whisper-ui.app/Contents/Resources/ffmpeg

whisper.cpp:
	git clone https://github.com/ggerganov/whisper.cpp.git
	cd whisper.cpp && sh ./models/download-ggml-model.sh medium
	cd whisper.cpp && make -j
	
ffmpeg:
	# note https://evermeet.cx/ffmpeg/apple-silicon-arm
	curl -Lso ffmpeg.zip https://evermeet.cx/ffmpeg/ffmpeg-117771-g07904231cb.zip
	gpg --import 0x1A660874.asc
	gpg --verify ffmpeg-117771-g07904231cb.zip.sig ffmpeg.zip
	unzip ffmpeg.zip
	xattr -dr com.apple.quarantine ffmpeg

clean:
	rm -rf whisper-ui whisper-ui.app

realclean: clean
	rm -rf whisper.cpp ffmpeg ffmpeg.zip whisper-cpp 