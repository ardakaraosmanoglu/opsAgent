package ai

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "net/http"
)

type Client struct {
    APIKey  string
    BaseURL string
    Model   string
}

func (c *Client) CreatePlan(ctx context.Context, userPrompt string, systemContext map[string]any) (*Plan, error) {
    if c.APIKey == "" || c.BaseURL == "" {
        return nil, errors.New("AI provider not configured")
    }

    payload := map[string]any{
        "model": c.Model,
        "messages": []map[string]string{
            {
                "role": "system",
                "content": `You are OpsAgent planner. You must only return safe structured JSON.
Never execute commands.
Never bypass approval.
Every command must be explicit.
Dangerous destructive commands must not be suggested.`,
            },
            {
                "role": "user",
                "content": userPrompt,
            },
        },
    }

    raw, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/chat/completions", bytes.NewReader(raw))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer "+c.APIKey)
    req.Header.Set("Content-Type", "application/json")

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode >= 300 {
        return nil, errors.New("AI provider returned non-success status")
    }

    var response struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }

    if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
        return nil, err
    }

    if len(response.Choices) == 0 {
        return nil, errors.New("no response from AI provider")
    }

    var plan Plan
    if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &plan); err != nil {
        return nil, err
    }

    return &plan, nil
}
