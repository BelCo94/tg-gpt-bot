package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiURL = "https://api.openai.com/v1/chat/completions"

type OpenAIConfig struct {
	Token        string `toml:"TOKEN"`
	Model        string `toml:"MODEL"`
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

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
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

	body, _ := ioutil.ReadAll(resp.Body)

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, err
	}

	if len(openAIResp.Choices) > 0 {
		return &openAIResp.Choices[0].Message, nil
	}
	return nil, fmt.Errorf("no message received from OpenAI")
}
