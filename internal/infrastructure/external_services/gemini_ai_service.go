package external_services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*

req = {
	content = [
		parts: [
			{Text: prompt}
		]
	]
}

res = {
	Candidates
}

*/
// -------------- gemini req & res dtos --------------
type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type RequestPayload struct {
	Contents []Content `json:"contents"`
}

type Candidate struct {
	Content struct {
		Parts []Part `json:"parts"`
	} `json:"content"`
}

type ResponsePayload struct {
	Candidates []Candidate `json:"candidates"`
}

// -------------- end of dto -------------------

type GeminiAIService struct {
	apiKey string
	model  string
	client *http.Client
}

func NewGeminiAIService(apiKey string) *GeminiAIService {
	defaultModel := "gemini-2.5-flash"
	return &GeminiAIService{
		apiKey: apiKey,
		model:  defaultModel,
		client: &http.Client{},
	}
}

func (as *GeminiAIService) GenerateContent(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", as.model, as.apiKey)

	payload := RequestPayload{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Gemini API payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed create gemini api request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// sending req to gemini api
	resp, err := as.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed create gemini api request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini return status code non-200: %v", resp.StatusCode)
	}

       var response ResponsePayload
       if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
	       return "", fmt.Errorf("failed to decode gemini response: %w", err)
       }
       if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
	       return response.Candidates[0].Content.Parts[0].Text, nil
       }
       return "", fmt.Errorf("gemini response has no content")
}
