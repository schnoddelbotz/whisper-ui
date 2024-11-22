build-darwin: clean ffmpeg whisper.cpp
	fyne version || go install fyne.io/fyne/v2/cmd/fyne@latest
	fyne package --release
	ln whisper.cpp/main whisper-ui.app/Contents/Resources/whisper-cpp
	ln ffmpeg whisper-ui.app/Contents/Resources/ffmpeg

# todo: zip-darwin, zip-linux

install-linux: whisper.cpp build-linux
	cp -v whisper-ui ~/bin
	# assuming deb/rpm ffmpeg is installed and ~/bin is in PATH
	cp -v whisper.cpp/main ~/bin/whisper-cpp
	mkdir -p ~/.whisper-ui
	cp -v whisper.cpp/models/ggml-medium.bin ~/.whisper-ui/ggml-medium.bin

build-linux:
	# todo: provide FyneApp.toml to create Linux .desktop file...?
	go build

whisper.cpp:
	git clone https://github.com/ggerganov/whisper.cpp.git
	cd whisper.cpp && make -j
	
ffmpeg:
	# note https://evermeet.cx/ffmpeg/apple-silicon-arm
	curl -Lso ffmpeg.zip https://evermeet.cx/ffmpeg/ffmpeg-117771-g07904231cb.zip
	test `md5 -q ffmpeg.zip` = af74364f54ccccd00faf37f67614cfb1 || exit 1
	unzip ffmpeg.zip
	xattr -dr com.apple.quarantine ffmpeg

clean:
	rm -rf whisper-ui whisper-ui.app

realclean: clean
	rm -rf whisper.cpp ffmpeg ffmpeg.zip
