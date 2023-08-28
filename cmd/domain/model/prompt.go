package model

type TranslationPrompt struct {
	JaTranslations []string `json:"ja_translations"`
	EnExplanation  string   `json:"en_explanation"`
}

type ExamplePrompt struct {
	Examples []*Example `json:"examples"`
}

type Example struct {
	Example     string `json:"example"`
	Translation string `json:"translation"`
}
