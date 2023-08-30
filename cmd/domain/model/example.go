package model

type Examples struct {
	Examples []*Example `json:"examples"`
}

type Example struct {
	Id          string `json:"id,omitempty"`
	Example     string `json:"example"`
	Translation string `json:"translation"`
}

func NewExample(id, example, translation string) *Example {
	return &Example{
		Id:          id,
		Example:     example,
		Translation: translation,
	}
}
