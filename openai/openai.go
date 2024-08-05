package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	textApiURL  = "https://api.openai.com/v1/chat/completions"
	imageApiURL = "https://api.openai.com/v1/images/generations"
)

type OpenAIConfig struct {
	Token        string `toml:"TOKEN"`
	Model        string `toml:"MODEL"`
	ImageModel   string `toml:"IMAGE_MODEL"`
	SystemPrompt string `toml:"SYSTEM_PROMPT"`
}

type MessageT struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type payloadT struct {
	Messages []MessageT `json:"messages"`
	Model    string     `json:"model"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message MessageT `json:"message"`
	} `json:"choices"`
}

type imagePayloadT struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model"`
	Size           string `json:"size"`
	NumberOfImages int    `json:"n"`
}

type OpenAIImageResponse struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}

func PreparePayload(messages []MessageT, config OpenAIConfig) ([]byte, error) {
	outputMessages := make([]MessageT, 0, len(messages)+1)
	systemElem := MessageT{
		Role:    "system",
		Content: config.SystemPrompt,
	}
	outputMessages = append(outputMessages, systemElem)
	for i := len(messages) - 1; i >= 0; i-- {
		outputMessages = append(outputMessages, messages[i])
	}
	data := payloadT{
		Model:    config.Model,
		Messages: outputMessages,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func AskOpenAI(payload []byte, config OpenAIConfig) (*MessageT, error) {

	req, err := http.NewRequest("POST", textApiURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, err
	}

	if len(openAIResp.Choices) > 0 {
		return &openAIResp.Choices[0].Message, nil
	}
	return nil, fmt.Errorf("no message received from OpenAI")
}

func PrepareImagePayload(prompt string, config OpenAIConfig) ([]byte, error) {
	data := imagePayloadT{
		Prompt:         prompt,
		Model:          config.ImageModel,
		Size:           "1024x1024",
		NumberOfImages: 1,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func AskOpenAIImage(payload []byte, config OpenAIConfig) (string, error) {
	req, err := http.NewRequest("POST", imageApiURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var openAIImageResp OpenAIImageResponse
	if err := json.Unmarshal(body, &openAIImageResp); err != nil {
		return "", err
	}

	if len(openAIImageResp.Data) > 0 {
		return openAIImageResp.Data[0].Url, nil
	}

	return "", fmt.Errorf("no image received from OpenAI")
}
