
build: whisper.cpp
	# building whisper-ui using fyne command line tool
	fyne version || go install fyne.io/tools/cmd/fyne@latest
	fyne package --release

build-darwin: clean ffmpeg build
	cp whisper.cpp/build/bin/whisper-cli \
		whisper.cpp/build/src/libwhisper.1.dylib whisper.cpp/build/ggml/src/libggml.0.dylib \
		whisper.cpp/build/ggml/src/libggml-cpu.0.dylib whisper.cpp/build/ggml/src/ggml-blas/libggml-blas.0.dylib \
		whisper.cpp/build/ggml/src/ggml-metal/libggml-metal.0.dylib whisper.cpp/build/ggml/src/libggml-base.0.dylib \
		whisper-ui.app/Contents/Resources/
	ln ffmpeg whisper-ui.app/Contents/Resources/ffmpeg

build-windows-on-darwin: ffmpeg-windows
	fyne version || go install fyne.io/tools/cmd/fyne@latest
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 fyne package -os windows -release

ffmpeg:
	# note https://evermeet.cx/ffmpeg/apple-silicon-arm
	curl -Lso ffmpeg.zip https://evermeet.cx/ffmpeg/ffmpeg-122342-g10b9984e59.zip
	test `md5 -q ffmpeg.zip` = 423b2a586493d39aab4e60c9a8545468 || exit 1
	unzip ffmpeg.zip
	xattr -dr com.apple.quarantine ffmpeg

ffmpeg-windows:
	curl -Lso ffmpeg-win.zip https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip
	# fixme ... sha256sum in GH workflow?
	# echo fa7d4d7e795db0e2503f49f105f46ed5852386f0cfdd819899be3b65ebde24fc ffmpeg-win.zip | sha256sum -c -
	7z x ffmpeg-win.zip


whisper.cpp:
	git clone https://github.com/ggerganov/whisper.cpp.git
	cd whisper.cpp && \
		git checkout 679bdb53dbcbfb3e42685f50c7ff367949fd4d48 && \
		patch -p1 < ../whisper.cpp.patch && \
		make build


zip-darwin:
	zip -r whisper-ui-macos-$(shell uname -m).zip whisper-ui.app

zip-windows:
	mkdir -p whisper-ui-windows64
	cp whisper-ui.exe whisper-ui-windows64/
	cp ffmpeg-8.0.1-essentials_build/bin/ffmpeg.exe whisper-ui-windows64/
	ls whisper.cpp/build
	ls whisper.cpp/build/Release
	cp whisper.cpp/build/bin/whisper-cli.exe whisper-ui-windows64/whisper-cli.exe
	cp whisper.cpp/build/bin/*.dll whisper-ui-windows64/
	7z a whisper-ui-windows64.zip whisper-ui-windows64

zip-linux: build
	rm -rf whisper-ui
	mkdir whisper-ui
	cd whisper-ui && tar -xJf ../whisper-ui.tar.xz
	# add whisper-cpp to fyne-generated archive
	ln whisper.cpp/build/bin/whisper-cli whisper-ui/usr/local/bin/whisper-cli
	# hack whisper-cpp into fyne-generated Makefile for (un)installation
	perl -pi -e 's@^install:@install:\n\tinstall -Dm00755 usr/local/bin/whisper-cli \$$(DESTDIR)\$$(PREFIX)/bin/whisper-cli@' whisper-ui/Makefile
	perl -pi -e 's@^uninstall:@uninstall:\n\t-rm \$$(DESTDIR)\$$(PREFIX)/bin/whisper-cli@' whisper-ui/Makefile
	perl -pi -e 's@^user-install:@user-install:\n\tinstall -Dm00755 usr/local/bin/whisper-cli \$$(DESTDIR)\$$(HOME)/.local/bin/whisper-cli@' whisper-ui/Makefile
	perl -pi -e 's@^user-uninstall:@user-uninstall:\n\t-rm \$$(DESTDIR)\$$(HOME)/.local/bin/whisper-cli@' whisper-ui/Makefile
	tar -czf whisper-ui-ubuntu-amd64.tar.gz whisper-ui/


clean:
	rm -rf whisper-ui whisper-ui.exe whisper-ui.app whisper-ui-*.zip whisper-ui-windows64 whisper-ui.tar.*

realclean: clean
	rm -rf whisper.cpp ffmpeg ffmpeg*.zip ffmpeg-7.1-essentials_build/
