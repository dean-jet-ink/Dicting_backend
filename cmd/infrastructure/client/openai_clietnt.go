package client

import (
	"context"
	"encoding/json"
	"english/cmd/domain/api"
	"english/cmd/domain/model"
	"english/config"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIAPI() api.ChatAIAPI {
	client := openai.NewClient(config.APIKey(config.OPENAI))
	return &OpenAIClient{
		client: client,
	}
}

func (c *OpenAIClient) GetTranslation(englishItem *model.EnglishItem) error {
	prompt := c.translationPrompt(englishItem.Content())

	resp, err := c.createChatCompletion(context.Background(), prompt)
	if err != nil {
		return err
	}

	answer := &model.Translation{}

	if err := json.Unmarshal([]byte(resp), answer); err != nil {
		return nil
	}

	englishItem.SetTranslations(answer.Translations)
	englishItem.SetEnExplanation(answer.EnExplanation)

	return nil
}

func (c *OpenAIClient) GetExample(englishItem *model.EnglishItem) error {
	prompt := c.examplePrompt(englishItem.Content())

	resp, err := c.createChatCompletion(context.Background(), prompt)
	if err != nil {
		return err
	}

	answer := &model.Examples{}

	if err := json.Unmarshal([]byte(resp), answer); err != nil {
		return err
	}

	englishItem.SetExamples(answer.Examples)

	return nil
}

func (c *OpenAIClient) createChatCompletion(ctx context.Context, prompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) translationPrompt(content string) string {
	answerExample := model.Translation{
		Translations: []string{
			"japanese",
			"japanese",
			"japanese",
		},
		EnExplanation: "english",
	}

	m, _ := json.Marshal(answerExample)

	prompt := fmt.Sprintf("[instructions]\ncreate three Japanese translations and a one-sentence basic English explanation of '%v' in JSON format.\n[answer of example]\n%v}", content, string(m))

	log.Printf("translationPrompt: %s\n", prompt)

	return prompt
}

func (c *OpenAIClient) examplePrompt(content string) string {
	answerExample := model.Examples{
		Examples: []*model.Example{
			{
				Example:     "english",
				Translation: "japanese",
			},
			{
				Example:     "english",
				Translation: "japanese",
			},
			{
				Example:     "english",
				Translation: "japanese",
			},
		},
	}

	m, _ := json.Marshal(answerExample)

	prompt := fmt.Sprintf("[instructions]\ncreate three sets of one-sentence basic English and its Japanese translation using %v in JSON format.\n[answer of example]\n%v}", content, string(m))

	log.Printf("examplePrompt: %s\n", prompt)

	return prompt
}
