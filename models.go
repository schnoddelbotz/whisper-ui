package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type model struct {
	name string
	size string
	sha  string
}

var availableModels = []model{
	{"tiny", "75 MiB", "bd577a113a864445d4c299885e0cb97d4ba92b5f"},
	{"tiny.en", "75 MiB ", "c78c86eb1a8faa21b369bcd33207cc90d64ae9df"},
	{"base", "142 MiB", "465707469ff3a37a2b9b8d8f89f2f99de7299dac"},
	{"base.en", "142 MiB", "137c40403d78fd54d454da0f9bd998f78703390c"},
	{"small", "466 MiB", "55356645c2b361a969dfd0ef2c5a50d530afd8d5"},
	{"small.en", "466 MiB", "db8a495a91d927739e50b3fc1cc4c6b8f6c2d022"},
	{"small.en-tdrz", "465 MiB", "b6c6e7e89af1a35c08e6de56b66ca6a02a2fdfa1"},
	{"medium", "1.5 GiB", "fd9727b6e1217c2f614f9b698455c4ffd82463b4"},
	{"medium.en", "1.5 GiB", "8c30f0e44ce9560643ebd10bbe50cd20eafd3723"},
	{"large-v1", "2.9 GiB", "b1caaf735c4cc1429223d5a74f0f4d0b9b59a299"},
	{"large-v2", "2.9 GiB", "0f4c8e34f21cf1a914c59d8b3ce882345ad349d6"},
	{"large-v2-q5_0", "1.1 GiB", "00e39f2196344e901b3a2bd5814807a769bd1630"},
	{"large-v3", "2.9 GiB", "ad82bf6a9043ceed055076d0fd39f5f186ff8062"},
	{"large-v3-q5_0", "1.1 GiB", "e6e2ed78495d403bef4b7cff42ef4aaadcfea8de"},
	{"large-v3-turbo", "1.5 GiB", "4af2b29d7ec73d781377bfd1758ca957a807e941"},
	{"large-v3-turbo-q5_0", "547 MiB", "e050f7970618a659205450ad97eb95a18d69c9ee"},
}

func getURLForModel(name string) string {
	// https://github.com/ggerganov/whisper.cpp/blob/master/models/download-ggml-model.sh
	src := "https://huggingface.co/ggerganov/whisper.cpp"
	pfx := "resolve/main/ggml"
	if strings.Contains(name, "tdrz") {
		src = "https://huggingface.co/akashmjn/tinydiarize-whisper.cpp"
		pfx = "resolve/main/ggml"
	}
	return fmt.Sprintf("%s/%s-%s.bin", src, pfx, name)
}

func getModel(name string) *model {
	for _, m := range availableModels {
		if name == m.name {
			return &m
		}
	}
	return nil
}

func getModels() []string {
	basedir := getModelsDir()
	models := []string{}
	modelsFullPath, _ := filepath.Glob(filepath.Join(basedir, "*.bin"))
	for _, m := range modelsFullPath {
		mod := strings.TrimPrefix(filepath.Base(m), "ggml-")
		mod = strings.TrimSuffix(mod, ".bin")
		models = append(models, mod)
	}
	return models
}

func modelIsInstalled(name string, models []string) bool {
	for _, m := range models {
		if m == name {
			return true
		}
	}
	return false
}
