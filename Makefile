build-darwin: clean ffmpeg build
	ln whisper.cpp/main whisper-ui.app/Contents/Resources/whisper-cpp
	ln ffmpeg whisper-ui.app/Contents/Resources/ffmpeg

build: whisper.cpp
	# todo: provide FyneApp.toml to create Linux .desktop file...?
	fyne version || go install fyne.io/fyne/v2/cmd/fyne@latest
	fyne package --release

build-windows-on-darwin: ffmpeg-windows
	fyne version || go install fyne.io/fyne/v2/cmd/fyne@latest
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 fyne package -os windows -release


ffmpeg:
	# note https://evermeet.cx/ffmpeg/apple-silicon-arm
	curl -Lso ffmpeg.zip https://evermeet.cx/ffmpeg/ffmpeg-117771-g07904231cb.zip
	test `md5 -q ffmpeg.zip` = af74364f54ccccd00faf37f67614cfb1 || exit 1
	unzip ffmpeg.zip
	xattr -dr com.apple.quarantine ffmpeg

ffmpeg-windows:
	curl -Lso ffmpeg-win.zip https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip
	# fixme ... sha256sum in GH workflow?
	# echo fa7d4d7e795db0e2503f49f105f46ed5852386f0cfdd819899be3b65ebde24fc ffmpeg-win.zip | sha256sum -c -
	7z x ffmpeg-win.zip


whisper.cpp:
	git clone https://github.com/ggerganov/whisper.cpp.git
	cd whisper.cpp && make -j main


zip-darwin:
	zip -r whisper-ui-macos-$(shell uname -m).zip whisper-ui.app

zip-windows:
	mkdir -p whisper-ui-windows64
	cp whisper-ui.exe whisper-ui-windows64/
	cp ffmpeg-7.1-essentials_build/bin/ffmpeg.exe whisper-ui-windows64/
	cp whisper.cpp/build/bin/Release/main.exe whisper-ui-windows64/whisper-cpp.exe
	cp whisper.cpp/build/bin/Release/*.dll whisper-ui-windows64/
	7z a whisper-ui-windows64.zip whisper-ui-windows64


clean:
	rm -rf whisper-ui whisper-ui.exe whisper-ui.app whisper-ui-*.zip whisper-ui-windows64

realclean: clean
	rm -rf whisper.cpp ffmpeg ffmpeg*.zip ffmpeg-7.1-essentials_build/
