package model

import "strings"

type EnglishItem struct {
	id             string
	content        string
	jaTranslations []string
	enExplanation  string
	examples       []*Example
	imgs           []*Img
	userId         string
}

func NewEnglishItem(id, content string, jaTranslations []string, enExplanation string, examples []*Example, imgs []*Img, userId string) *EnglishItem {
	return &EnglishItem{
		id:             id,
		content:        content,
		jaTranslations: jaTranslations,
		enExplanation:  enExplanation,
		examples:       examples,
		imgs:           imgs,
		userId:         userId,
	}
}

func (e *EnglishItem) JoinJaTranslations() string {
	return strings.Join(e.jaTranslations, ",")
}

func (e *EnglishItem) Id() string {
	return e.id
}

func (e *EnglishItem) Content() string {
	return e.content
}

func (e *EnglishItem) JaTranslations() []string {
	return e.jaTranslations
}

func (e *EnglishItem) EnExplanation() string {
	return e.enExplanation
}

func (e *EnglishItem) Examples() []*Example {
	return e.examples
}

func (e *EnglishItem) Imgs() []*Img {
	return e.imgs
}

func (e *EnglishItem) UserId() string {
	return e.userId
}

func (e *EnglishItem) SetId(id string) {
	e.id = id
}

func (e *EnglishItem) SetContent(content string) {
	e.content = content
}

func (e *EnglishItem) SetJaTranslations(jaTranslations []string) {
	e.jaTranslations = jaTranslations
}

func (e *EnglishItem) SetJaTranslationsFromStr(jaTranslations string) {
	e.jaTranslations = strings.Split(jaTranslations, ",")
}

func (e *EnglishItem) SetEnExplanation(enExplanation string) {
	e.enExplanation = enExplanation
}

func (e *EnglishItem) SetExamples(examples []*Example) {
	e.examples = examples
}

func (e *EnglishItem) SetImgURLs(imgs []*Img) {
	e.imgs = imgs
}

func (e *EnglishItem) SetUserId(userId string) {
	e.userId = userId
}
