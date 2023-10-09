package validator

import (
	"english/cmd/domain/model"
	"english/myerror"
	"errors"
	"fmt"
)

type EnglishItemValidator interface {
	EnglishItemValidate(englishItem *model.EnglishItem) error
}

type englishItemValidatorImpl struct {
}

func NewEnglishItemValidator() EnglishItemValidator {
	return &englishItemValidatorImpl{}
}

var proficiencies = map[string]struct{}{
	"Learning":   {},
	"Understand": {},
	"Mastered":   {},
}

func (v *englishItemValidatorImpl) EnglishItemValidate(englishItem *model.EnglishItem) error {
	proficiency := englishItem.Proficiency()

	if _, ok := proficiencies[proficiency]; !ok {
		return fmt.Errorf("%v: %w", myerror.ErrValidation, errors.New("unexpected proficiency"))
	}

	return nil
}
