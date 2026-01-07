package main

import (
	"log"
	"os"
	"sort"
)

type resources struct {
	ffmpeg     string
	whispercpp string
	tmpfile    string
}

var languageMap = map[string]string{
	// curl -s https://raw.githubusercontent.com/ggml-org/whisper.cpp/refs/heads/master/src/whisper.cpp | \
	//   grep -A100 'g_lang = {' | tail -n +2 | \
	//   perl -pe 's@.*"([^"]+)".*"([^"]+)".*@"\"".ucfirst($2)."\": \"$1\",";@e' | sort
	"Afrikaans":      "af",
	"Albanian":       "sq",
	"Amharic":        "am",
	"Arabic":         "ar",
	"Armenian":       "hy",
	"Assamese":       "as",
	"Azerbaijani":    "az",
	"Bashkir":        "ba",
	"Basque":         "eu",
	"Belarusian":     "be",
	"Bengali":        "bn",
	"Bosnian":        "bs",
	"Breton":         "br",
	"Bulgarian":      "bg",
	"Cantonese":      "yue",
	"Catalan":        "ca",
	"Chinese":        "zh",
	"Croatian":       "hr",
	"Czech":          "cs",
	"Danish":         "da",
	"Dutch":          "nl",
	"English":        "en",
	"Estonian":       "et",
	"Faroese":        "fo",
	"Finnish":        "fi",
	"French":         "fr",
	"Galician":       "gl",
	"Georgian":       "ka",
	"German":         "de",
	"Greek":          "el",
	"Gujarati":       "gu",
	"Haitian creole": "ht",
	"Hausa":          "ha",
	"Hawaiian":       "haw",
	"Hebrew":         "he",
	"Hindi":          "hi",
	"Hungarian":      "hu",
	"Icelandic":      "is",
	"Indonesian":     "id",
	"Italian":        "it",
	"Japanese":       "ja",
	"Javanese":       "jw",
	"Kannada":        "kn",
	"Kazakh":         "kk",
	"Khmer":          "km",
	"Korean":         "ko",
	"Lao":            "lo",
	"Latin":          "la",
	"Latvian":        "lv",
	"Lingala":        "ln",
	"Lithuanian":     "lt",
	"Luxembourgish":  "lb",
	"Macedonian":     "mk",
	"Malagasy":       "mg",
	"Malay":          "ms",
	"Malayalam":      "ml",
	"Maltese":        "mt",
	"Maori":          "mi",
	"Marathi":        "mr",
	"Mongolian":      "mn",
	"Myanmar":        "my",
	"Nepali":         "ne",
	"Norwegian":      "no",
	"Nynorsk":        "nn",
	"Occitan":        "oc",
	"Pashto":         "ps",
	"Persian":        "fa",
	"Polish":         "pl",
	"Portuguese":     "pt",
	"Punjabi":        "pa",
	"Romanian":       "ro",
	"Russian":        "ru",
	"Sanskrit":       "sa",
	"Serbian":        "sr",
	"Shona":          "sn",
	"Sindhi":         "sd",
	"Sinhala":        "si",
	"Slovak":         "sk",
	"Slovenian":      "sl",
	"Somali":         "so",
	"Spanish":        "es",
	"Sundanese":      "su",
	"Swahili":        "sw",
	"Swedish":        "sv",
	"Tagalog":        "tl",
	"Tajik":          "tg",
	"Tamil":          "ta",
	"Tatar":          "tt",
	"Telugu":         "te",
	"Thai":           "th",
	"Tibetan":        "bo",
	"Turkish":        "tr",
	"Turkmen":        "tk",
	"Ukrainian":      "uk",
	"Urdu":           "ur",
	"Uzbek":          "uz",
	"Vietnamese":     "vi",
	"Welsh":          "cy",
	"Yiddish":        "yi",
	"Yoruba":         "yo",
}

func getTempFileName() string {
	// just interested in a name in $TMP that does not exist.
	// CreateTemp ensures this ... is there a way to avoid creation to just get name?
	f, err := os.CreateTemp("", "whisper-ui-*.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return f.Name()
}

func languages() []string {
	languages := make([]string, 0, len(languageMap))
	for language := range languageMap {
		languages = append(languages, language)
	}
	sort.Strings(languages)
	return languages
}
