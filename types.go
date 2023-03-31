package main

type ChatRequest struct {
	Model            string              `json:"model"`
	Messages         []map[string]string `json:"messages"`
	Temperature      float64             `json:"temperature"`
	TopP             float64             `json:"top_p"`
	N                int                 `json:"n"`
	PresencePenalty  float64             `json:"presence_penalty"`
	FrequencyPenalty float64             `json:"frequency_penalty"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type envs struct {
	OpenAIKey   string
	AiImgSecret string
	AiImgKey    string
}
