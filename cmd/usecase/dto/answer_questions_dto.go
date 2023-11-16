package dto

type Answer struct {
	Index    uint   `json:"index"`
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer" validate:"required"`
}

type AnswerQuestionsInput struct {
	EnglishItemId string    `json:"english_item_id" validate:"required"`
	Content       string    `json:"content" validate:"required"`
	Answers       []*Answer `json:"answers" validate:"dive"`
}

type Advice struct {
	Index  uint   `json:"index"`
	Advice string `json:"advice"`
}

type AnswerQuestionsOutput struct {
	AdviceList []*Advice `json:"advice_list"`
}
