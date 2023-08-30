package dto

type ProposalEnglishItemRequest struct {
	Content string `form:"content" validate:"required,lt=255"`
}

type ProposalEnglishItemResponse struct {
	Content        string     `json:"content"`
	JaTranslations []string   `json:"ja_translations,omitempty"`
	EnExplanation  string     `json:"en_explanation,omitempty"`
	Examples       []*Example `json:"examples,omitempty"`
}
