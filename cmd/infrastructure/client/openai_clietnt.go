package client

import (
	"context"
	"encoding/json"
	"english/cmd/domain/api"
	"english/cmd/domain/model"
	"english/config"
	"fmt"
	"log"
	"sync"

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

	res, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) GetTranslations(englishItem *model.EnglishItem) error {
	prompt := c.translationsPrompt(englishItem.Content())

	res, err := c.createChatCompletion(context.Background(), prompt)
	if err != nil {
		return err
	}

	answer := &model.Translation{}

	if err := json.Unmarshal([]byte(res), answer); err != nil {
		return nil
	}

	englishItem.SetTranslations(answer.Translations)
	englishItem.SetEnExplanation(answer.EnExplanation)

	return nil
}

func (c *OpenAIClient) GetExamples(englishItem *model.EnglishItem) error {
	prompt := c.examplesPrompt(englishItem.Content())

	res, err := c.createChatCompletion(context.Background(), prompt)
	if err != nil {
		return err
	}

	answer := &model.Examples{}

	if err := json.Unmarshal([]byte(res), answer); err != nil {
		return err
	}

	englishItem.SetExamples(answer.Examples)

	return nil
}

func (c *OpenAIClient) GetTranslation(content string) (string, error) {
	prompt := c.translationPrompt(content)

	res, err := c.createChatCompletion(context.Background(), prompt)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (c *OpenAIClient) GetExplanation(content string) (string, error) {
	prompt := c.explanationPrompt(content)

	res, err := c.createChatCompletion(context.Background(), prompt)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (c *OpenAIClient) GetExample(content string) (*model.Example, error) {
	prompt := c.examplePrompt(content)

	res, err := c.createChatCompletion(context.Background(), prompt)

	log.Println(res)

	if err != nil {
		return nil, err
	}

	answer := &model.Example{}

	if err = json.Unmarshal([]byte(res), answer); err != nil {
		return nil, err
	}

	return answer, nil
}

func (c *OpenAIClient) translationsPrompt(content string) string {
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

	log.Printf("translationsPrompt: %s\n", prompt)

	return prompt
}

func (c *OpenAIClient) translationPrompt(content string) string {
	prompt := fmt.Sprintf("[instructions]\nTranslate '%v' into Japanese.\n[format]Only Japanese translation. Do not output English.", content)

	return prompt
}

func (c *OpenAIClient) explanationPrompt(content string) string {
	prompt := fmt.Sprintf("[instructions]\ncreate a one-sentence basic English explanation of '%v'.}", content)

	return prompt
}

func (c *OpenAIClient) examplesPrompt(content string) string {
	answerExample := model.Examples{
		Examples: []*model.Example{
			{
				Example:     "",
				Translation: "",
			},
			{
				Example:     "",
				Translation: "",
			},
			{
				Example:     "",
				Translation: "",
			},
		},
	}

	m, _ := json.Marshal(answerExample)

	prompt := fmt.Sprintf("[instructions]\ncreate three sets of one-sentence basic English and its Japanese translation using '%v' in JSON format.\n[json format]\n%v}", content, string(m))

	return prompt
}

func (c *OpenAIClient) examplePrompt(content string) string {
	answerExample := model.Example{
		Example:     "",
		Translation: "",
	}

	m, _ := json.Marshal(answerExample)

	prompt := fmt.Sprintf("[instructions]\ncreate a one-sentence basic English and its Japanese translation using '%v' in JSON format.\n[json format]\n%v}", content, string(m))

	return prompt
}

func (c *OpenAIClient) GetQuestion(content string) (string, error) {
	prompt := c.questionPrompt(content)

	res, err := c.createChatCompletion(context.Background(), prompt)
	if err != nil {
		return "", err
	}

	log.Println(res)

	return res, nil
}

func (c *OpenAIClient) questionPrompt(content string) string {
	return fmt.Sprintf("[instructions]\nProvide a Japanese sentence that includes the phrase '%v' in a natural context. However, exclude expressions unique to Japanese.\n[format]\nOnly Japanese sentence. Do not output English.", content)
}

func (c *OpenAIClient) GetAdvice(answers []*model.Output) error {
	errChan := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	for _, answer := range answers {
		wg.Add(1)

		go func(answer *model.Output) {
			defer wg.Done()

			prompt := c.advicePrompt(answer.Content(), answer.Question(), answer.Answer())

			res, err := c.createChatCompletion(ctx, prompt)
			if err != nil {
				errChan <- err
				return
			}

			answer.SetAdvice(res)
		}(answer)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		cancel()
		return err
	}

	return nil
}

func (c *OpenAIClient) advicePrompt(keyword, question, answer string) string {
	return fmt.Sprintf("'%v'を使った英作文をしています。以下の英訳についてアドバイスをください。日本語でお願いします。\n本文: %v\n英訳: %v", keyword, question, answer)
}
