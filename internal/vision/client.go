package vision

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	APIKey string
	Model  string
}

func NewClient() *Client {
	key := os.Getenv("ANTHROPIC_API_KEY")
	return &Client{APIKey: key, Model: "claude-opus-4-6"}
}

func (c *Client) Analyze(ctx context.Context, imagePath, prompt string) (string, error) {
	if c.APIKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY not set")
	}
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("reading image: %w", err)
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	mediaType := "image/png"
	if len(imagePath) > 4 {
		ext := imagePath[len(imagePath)-4:]
		if ext == ".jpg" || ext == "jpeg" {
			mediaType = "image/jpeg"
		}
		if ext == ".gif" {
			mediaType = "image/gif"
		}
		if ext == ".webp" {
			mediaType = "image/webp"
		}
	}
	body := map[string]any{
		"model":      c.Model,
		"max_tokens": 1024,
		"messages": []any{map[string]any{
			"role": "user",
			"content": []any{
				map[string]any{"type": "image", "source": map[string]any{"type": "base64", "media_type": mediaType, "data": b64}},
				map[string]any{"type": "text", "text": prompt},
			},
		}},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(bodyBytes))
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	json.Unmarshal(respBody, &result)
	if result.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}
	if len(result.Content) > 0 {
		return result.Content[0].Text, nil
	}
	return "", fmt.Errorf("empty response")
}
