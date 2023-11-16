package dto

type GetQuestionInput struct {
	Content string `form:"content" validate:"required"`
}

type GetQuestionOutput struct {
	Question string `json:"question"`
}
