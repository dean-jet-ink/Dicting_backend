package usecase

import (
	"english/cmd/domain/api"
	"english/cmd/domain/model"
	"english/cmd/usecase/dto"
)

type ProposalEnglishItemUsecase interface {
	// Proposal(req *dto.ProposalEnglishItemRequest) *dto.ProposalEnglishItemResponse
	Proposal(req *dto.ProposalEnglishItemRequest) (*dto.ProposalEnglishItemResponse, error)
}

type ProposalEnglishItemUsecaseImpl struct {
	chatAIAPI api.ChatAIAPI
}

func NewProposalEnglishItemUsecase(chatAIAPI api.ChatAIAPI) ProposalEnglishItemUsecase {
	return &ProposalEnglishItemUsecaseImpl{
		chatAIAPI: chatAIAPI,
	}
}

func (pu *ProposalEnglishItemUsecaseImpl) Proposal(req *dto.ProposalEnglishItemRequest) (*dto.ProposalEnglishItemResponse, error) {
	englishItem := model.NewEnglishItem("", req.Content, nil, "", nil, nil, "")

	// goroutineで記載
	if err := pu.chatAIAPI.GetTranslation(englishItem); err != nil {
		return nil, err
	}

	if err := pu.chatAIAPI.GetExample(englishItem); err != nil {
		return nil, err
	}

	examples := []*dto.Example{}
	for _, example := range englishItem.Examples() {
		exampleDto := &dto.Example{
			Example:     example.Example,
			Translation: example.Translation,
		}
		examples = append(examples, exampleDto)
	}

	resp := &dto.ProposalEnglishItemResponse{
		Content:        englishItem.Content(),
		JaTranslations: englishItem.JaTranslations(),
		EnExplanation:  englishItem.EnExplanation(),
		Examples:       examples,
	}

	return resp, nil
}
