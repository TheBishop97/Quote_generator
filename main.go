package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "io/ioutil"
    "bytes"
)

const apiKey = "YOUR_OPENAI_API_KEY"

type ChatGPTResponse struct {
    Choices []struct {
        Text string `json:"text"`
    } `json:"choices"`
}

func fetchQuote() (string, error) {
    url := "https://api.openai.com/v1/engines/davinci-codex/completions"
    prompt := `Generate an inspirational quote.`
    data := map[string]interface{}{
        "prompt": prompt,
        "max_tokens": 60,
    }
    jsonData, _ := json.Marshal(data)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var gptResponse ChatGPTResponse
    err = json.Unmarshal(body, &gptResponse)
    if err != nil {
        return "", err
    }

    if len(gptResponse.Choices) > 0 {
        return gptResponse.Choices[0].Text, nil
    }

    return "No quote available", nil
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
    quote, err := fetchQuote()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    jsonResponse, err := json.Marshal(map[string]string{"quote": quote})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}

func main() {
    http.HandleFunc("/quote", quoteHandler)
    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", nil)
}
