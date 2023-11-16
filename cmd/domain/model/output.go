package model

import "time"

type Output struct {
	id            string
	englishItemId string
	content       string
	question      string
	answer        string
	advice        string
	index         uint
	createdAt     time.Time
}

func NewOutput(id, englishItemId, content, question, answer, advice string, index uint, createdAt time.Time) *Output {
	return &Output{
		id:            id,
		englishItemId: englishItemId,
		content:       content,
		question:      question,
		answer:        answer,
		advice:        advice,
		index:         index,
		createdAt:     createdAt,
	}
}

func (o *Output) ID() string {
	return o.id
}

func (o *Output) Content() string {
	return o.content
}

func (o *Output) EnglishItemId() string {
	return o.englishItemId
}

func (o *Output) Question() string {
	return o.question
}

func (o *Output) Answer() string {
	return o.answer
}

func (o *Output) Advice() string {
	return o.advice
}

func (o *Output) Index() uint {
	return o.index
}

func (o *Output) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Output) SetID(id string) {
	o.id = id
}

func (o *Output) SetEnglishItemId(englishItemId string) {
	o.englishItemId = englishItemId
}

func (o *Output) SetQuestion(question string) {
	o.question = question
}

func (o *Output) SetAnswer(answer string) {
	o.answer = answer
}

func (o *Output) SetAdvice(advice string) {
	o.advice = advice
}

func (o *Output) SetIndex(index uint) {
	o.index = index
}

func (o *Output) SetCreatedAt(createdAt time.Time) {
	o.createdAt = createdAt
}
