package model

import "strings"

type EnglishItem struct {
	id            string
	content       string
	translations  []string
	enExplanation string
	examples      []*Example
	imgs          []*Img
	userId        string
}

func NewEnglishItem(id, content string, translations []string, enExplanation string, examples []*Example, imgs []*Img, userId string) *EnglishItem {
	return &EnglishItem{
		id:            id,
		content:       content,
		translations:  translations,
		enExplanation: enExplanation,
		examples:      examples,
		imgs:          imgs,
		userId:        userId,
	}
}

func (e *EnglishItem) JoinTranslations() string {
	return strings.Join(e.translations, ",")
}

func (e *EnglishItem) Id() string {
	return e.id
}

func (e *EnglishItem) Content() string {
	return e.content
}

func (e *EnglishItem) Translations() []string {
	return e.translations
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

func (e *EnglishItem) SetTranslations(translations []string) {
	e.translations = translations
}

func (e *EnglishItem) SetTranslationsFromStr(translations string) {
	e.translations = strings.Split(translations, ",")
}

func (e *EnglishItem) SetEnExplanation(enExplanation string) {
	e.enExplanation = enExplanation
}

func (e *EnglishItem) SetExamples(examples []*Example) {
	e.examples = examples
}

func (e *EnglishItem) SetImgs(imgs []*Img) {
	e.imgs = imgs
}

func (e *EnglishItem) SetUserId(userId string) {
	e.userId = userId
}
