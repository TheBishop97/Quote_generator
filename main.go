package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// ChatGPTResponse is a struct that represents the response from the ChatGPT API.
// It contains a slice of Choices, where each Choice represents a possible completion of the prompt.
type ChatGPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func fetchQuote() (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	url := "https://api.openai.com/v1/chat/completions" // Updated endpoint for chat model
	prompt := `Generate an inspirational quote.`
	data := map[string]interface{}{
		"model": "gpt-3.5-turbo", // Updated model
		"messages": []map[string]string{
			{"role": "system", "content": "You are an inspirational quote generator."},
			{"role": "user", "content": prompt},
		},
	}
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	fmt.Println("Request URL:", url)
	fmt.Println("Request Headers:", req.Header)
	fmt.Println("Request Body:", string(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Response Body:", string(body))

	var gptResponse ChatGPTResponse
	err = json.Unmarshal(body, &gptResponse)
	if err != nil {
		return "", err
	}

	if len(gptResponse.Choices) > 0 {
		return gptResponse.Choices[0].Message.Content, nil
	}

	return "No quote available", nil
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	quote, err := fetchQuote()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Quote:", quote)

	jsonResponse, err := json.Marshal(map[string]string{"quote": quote})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	http.HandleFunc("/quote", quoteHandler)
	fmt.Println("Open this URL in your browser: http://localhost:8080/quote")
	http.ListenAndServe(":8080", nil)
}
