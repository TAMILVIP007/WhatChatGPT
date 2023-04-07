package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"io"
	"net/http"
)

// GenAIimg function is used to generate anime images based on the input image bytes
func GenAIimg(imageBytes []byte) ([]byte, error) {
	// Encode image in base64 string
	encodedString := base64.StdEncoding.EncodeToString(imageBytes)
	// Create payload with all the required info into a map
	payload := map[string]interface{}{
		"parameter": map[string]string{
			"rsp_media_type": "jpg",
		},
		"extra": map[string]interface{}{},
		"media_info_list": []map[string]interface{}{
			{
				"media_data": encodedString,
				"media_profiles": map[string]string{
					"media_data_type": "jpg",
				},
				"media_extra": map[string]interface{}{},
			},
		},
	}
	// Marshal payload map into json bytes
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Couldn't marshal json data: %v", err.Error())
	}
	// Initialize http client to send requests
	client := &http.Client{}
	// Create a http POST request and set headers
	req, err := http.NewRequest("POST", fmt.Sprintf("https://openapi.mtlab.meitu.com/v1/stable_diffusion_anime?api_key=%s&api_secret=%s", env.AiImgKey, env.AiImgSecret), nil)
	if err != nil {
		log.Printf("Can't Send http req: %v", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	// Set request body to the payload we created earlier
	req.Body = io.NopCloser(bytes.NewReader(payloadBytes))
	// Send the request and set the response to res variable
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Close the response body when the function returns
	defer resp.Body.Close()

	// Read all bytes from the response body
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Unmarshal response to a map
	var respMap map[string]interface{}
	err = json.Unmarshal(respBytes, &respMap)
	if err != nil {
		return nil, err
	}

	// Get media data from the repose and convert it back to bytes
	mediaList, ok := respMap["media_info_list"].([]interface{})
	if !ok || len(mediaList) == 0 {
		return nil, fmt.Errorf("invalid response format: media_info_list field is missing or has an invalid format")
	}
	mediaData, ok := mediaList[0].(map[string]interface{})["media_data"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid response format: media_data field is missing or has an invalid format")
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(mediaData)
	if err != nil {
		return nil, err
	}
	return decodedBytes, nil
}

// GetAiTextResponse function is used to get AI text response for the given input message
func GetAiTextResponse(msg string) (string, error) {
	// Define URL and payload structure
	url := "https://api.openai.com/v1/chat/completions"
	payload := ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []map[string]string{
			{
				"role":    "user",
				"content": msg,
			},
		},
		Temperature:      1.0,
		TopP:             1.0,
		N:                1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	}

	// Marshall structure into JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling payload: %v", err.Error())
	}

	// Create new http post request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err.Error())
	}
	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", env.OpenAIKey))

	// Initialize the client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err.Error())
	}

	// Close response body after it is read
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err.Error())
	}

	// Unmarshal the chat response from bytes
	var chatResponse ChatResponse	
	if err = json.Unmarshal(responseBody, &chatResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling response body: %v", err.Error())
	}
	if len(chatResponse.Choices) == 0 {
		log.Printf("Resp From OPenAi: %v", err.Error())
		return "", fmt.Errorf("error getting response from OpenAI")
	}
	content := chatResponse.Choices[0].Message.Content
	return content, nil
}
