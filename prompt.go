package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []Choice `json:"choices"`
}

func HandlePrompt(c echo.Context) error {
	var promptReq PromptRequest
	if err := c.Bind(&promptReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Prepare request body for Groq API
	reqBody := RequestBody{
		Model: os.Getenv("MODEL_NAME"),
		Messages: []Message{
			{
				Role:    "user",
				Content: promptReq.Prompt,
			},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create request"})
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send request to Groq"})
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var groqResp ResponseBody
	json.Unmarshal(body, &groqResp)

	// Return AI response
	if len(groqResp.Choices) > 0 {
		return c.JSON(http.StatusOK, map[string]string{"response": groqResp.Choices[0].Message.Content})
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "No response from Groq"})
	}
}
