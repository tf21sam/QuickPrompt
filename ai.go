package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type GroqRequest struct {
	Model   string `json:"model"`
	Message []struct {
		Role    string `jsn:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GetAIResponse(pormpt string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	model := os.Getenv("MODEL_NAME")

	reqBody := map[string]interface{}{
		"model": model,
		"message": []map[string]string{
			{"role": "user", "content": pormpt},
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var groqResp GroqResponse
	if err := json.Unmarshal(respBody, &groqResp); err != nil {
		return "", err
	}
	return groqResp.Choices[0].Message.Content, nil

}
