package model

type Translation struct {
	JaTranslations []string `json:"ja_translations"`
	EnExplanation  string   `json:"en_explanation"`
}
