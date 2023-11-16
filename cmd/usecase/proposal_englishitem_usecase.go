package usecase

import (
	"english/cmd/domain/api"
	"english/cmd/domain/model"
	"english/cmd/usecase/dto"
	"sync"
)

type ProposalEnglishItemUsecase interface {
	Proposal(req *dto.ProposalEnglishItemRequest) (*dto.ProposalEnglishItemResponse, error)
	ProposalTranslation(req *dto.ProposalEnglishItemRequest) (string, error)
	ProposalExplanation(req *dto.ProposalEnglishItemRequest) (string, error)
	ProposalExample(req *dto.ProposalEnglishItemRequest) (*dto.Example, error)
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
	englishItem := model.NewEnglishItem("", req.Content, nil, "", nil, nil, "", model.Learning, 0)

	wg := sync.WaitGroup{}
	errChan := make(chan error)

	wg.Add(2)
	go func() {
		if err := pu.chatAIAPI.GetTranslations(englishItem); err != nil {
			errChan <- err
		}
		wg.Done()
	}()

	go func() {
		if err := pu.chatAIAPI.GetExamples(englishItem); err != nil {
			errChan <- err
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return nil, err
		}
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
		JaTranslations: englishItem.Translations(),
		EnExplanation:  englishItem.EnExplanation(),
		Examples:       examples,
	}

	return resp, nil
}

func (pu *ProposalEnglishItemUsecaseImpl) ProposalTranslation(req *dto.ProposalEnglishItemRequest) (string, error) {
	translation, err := pu.chatAIAPI.GetTranslation(req.Content)
	if err != nil {
		return "", err
	}

	return translation, nil
}

func (pu *ProposalEnglishItemUsecaseImpl) ProposalExplanation(req *dto.ProposalEnglishItemRequest) (string, error) {
	explanation, err := pu.chatAIAPI.GetExplanation(req.Content)
	if err != nil {
		return "", err
	}

	return explanation, nil
}

func (pu *ProposalEnglishItemUsecaseImpl) ProposalExample(req *dto.ProposalEnglishItemRequest) (*dto.Example, error) {
	example, err := pu.chatAIAPI.GetExample(req.Content)
	if err != nil {
		return nil, err
	}

	exampleDTO := &dto.Example{
		Example:     example.Example,
		Translation: example.Translation,
	}

	return exampleDTO, nil
}
