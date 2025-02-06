package llm

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
)

// Message represents a chat message
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// Usage tracks token usage from API response
type Usage struct {
    PromptTokens     int `json:"prompt_tokens"`
    CompletionTokens int `json:"completion_tokens"`
    TotalTokens      int `json:"total_tokens"`
}

// CompletionResponse represents API response structure
type CompletionResponse struct {
    Choices []struct {
        Delta struct {
            Content string `json:"content"`
        } `json:"delta"`
        FinishReason string `json:"finish_reason,omitempty"`
    } `json:"choices"`
    Usage *Usage `json:"usage"`
}

// Client handles interactions with the OpenAI API
type Client struct {
    config     *Config
    httpClient *http.Client
    baseURL    string
    model      string
}

func New(config *Config) (*Client, error) {
    if err := config.Validate(); err != nil {
        return nil, err
    }
    return &Client{
        config:     config,
        httpClient: &http.Client{},
        baseURL:    "https://api.openai.com/v1/chat/completions",
        model:      config.Model,
    }, nil
}

func (c *Client) Chat(ctx context.Context, messages []Message) (chan string, chan uint64, chan error) {
    requestStart := time.Now()
    responses := make(chan string)
    tokens := make(chan uint64)
    errors := make(chan error)

    var debug *log.Logger
    if os.Getenv("DEBUG") != "" {
        debug = log.New(os.Stderr, "\n[LLM] ", log.Ltime|log.Lmicroseconds)
        debug.Printf("Starting request with model=%s msgCount=%d", c.model, len(messages))
        os.Stderr.Sync()
    }

    go func() {
        defer close(responses)
        defer close(tokens)
        defer close(errors)

        // First, make a non-streaming request to get token count
        nonStreamBody, err := json.Marshal(map[string]interface{}{
            "model":    c.model,
            "messages": messages,
        })
        if err != nil {
            if debug != nil {
                debug.Printf("Error marshaling token request: %v", err)
            }
            errors <- fmt.Errorf("marshaling token request: %w", err)
            return
        }

        if debug != nil {
            debug.Printf("Non-streaming request body: %s", string(nonStreamBody))
        }

        nonStreamReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, strings.NewReader(string(nonStreamBody)))
        if err != nil {
            if debug != nil {
                debug.Printf("Error creating token request: %v", err)
            }
            errors <- fmt.Errorf("creating token request: %w", err)
            return
        }
        nonStreamReq.Header.Set("Content-Type", "application/json")
        nonStreamReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
        
        nonStreamResp, err := c.httpClient.Do(nonStreamReq)
        if err != nil {
            if debug != nil {
                debug.Printf("Error making non-streaming request: %v", err)
            }
            errors <- fmt.Errorf("token count request failed: %w", err)
            return
        }
        defer nonStreamResp.Body.Close()
        
        if debug != nil {
            debug.Printf("Non-streaming response status: %s", nonStreamResp.Status)
        }
        body, err := io.ReadAll(nonStreamResp.Body)
        if err != nil {
            if debug != nil {
                debug.Printf("Error reading response body: %v", err)
            }
            errors <- fmt.Errorf("reading token count response: %w", err)
            return
        }
        if debug != nil {
            debug.Printf("Non-streaming response body: %s", string(body))
        }

        if nonStreamResp.StatusCode != http.StatusOK {
            if debug != nil {
                debug.Printf("Error response from API: %s", nonStreamResp.Status)
            }
            errors <- fmt.Errorf("API error: %s", nonStreamResp.Status)
            return
        }

        var result CompletionResponse
        if err := json.Unmarshal(body, &result); err != nil {
            if debug != nil {
                debug.Printf("Error parsing non-streaming response: %v", err)
            }
            errors <- fmt.Errorf("parsing token count response: %w", err)
            return
        }

        if result.Usage == nil {
            if debug != nil {
                debug.Printf("No usage info in response")
            }
            errors <- fmt.Errorf("no token usage in response")
            return
        }
        // GPT-3.5 has 4k context window
        const maxContextTokens = 4096
        usedPct := float64(result.Usage.TotalTokens) / float64(maxContextTokens) * 100
        
        if os.Getenv("DEBUG") != "" {
            debug.Printf("Tokens Used: %d total\n" +
                "  %d input + %d output\n" +
                "  %.0f%% of 4k limit\n" +
                "  %.1fs",
                result.Usage.TotalTokens,
                result.Usage.PromptTokens,
                result.Usage.CompletionTokens,
                usedPct,
                time.Since(requestStart).Seconds())
        }

        totalTokens := uint64(result.Usage.TotalTokens)

        // Send token count BEFORE starting streaming request
        totalTokens = uint64(result.Usage.TotalTokens)
        if debug != nil {
            debug.Printf("Sending token count: %d", totalTokens)
        }
        select {
        case tokens <- totalTokens:
            if debug != nil {
                debug.Printf("Successfully sent token count: %d", totalTokens)
            }
        case <-ctx.Done():
            if debug != nil {
                debug.Printf("Context cancelled before sending token count")
            }
            return
        }

        // Now do streaming request
        reqBody, err := json.Marshal(map[string]interface{}{
            "model":    c.model,
            "messages": messages,
            "stream":   true,
        })
        if err != nil {
            if debug != nil {
                debug.Printf("Error marshaling streaming request: %v", err)
            }
            errors <- fmt.Errorf("failed to marshal request: %w", err)
            return
        }

        if debug != nil {
            debug.Printf("Streaming request body: %s", string(reqBody))
        }

        req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, strings.NewReader(string(reqBody)))
        if err != nil {
            if debug != nil {
                debug.Printf("Error creating streaming request: %v", err)
            }
            errors <- fmt.Errorf("failed to create request: %w", err)
            return
        }

        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

        resp, err := c.httpClient.Do(req)
        if err != nil {
            if debug != nil {
                debug.Printf("Error sending streaming request: %v", err)
            }
            errors <- fmt.Errorf("failed to send request: %w", err)
            return
        }
        defer resp.Body.Close()

        if debug != nil {
            debug.Printf("Streaming response status: %s", resp.Status)
        }

        if resp.StatusCode != http.StatusOK {
            body, _ := io.ReadAll(resp.Body)
            if debug != nil {
                debug.Printf("Streaming error response: %s", string(body))
            }
            errors <- fmt.Errorf("API error: %s: %s", resp.Status, string(body))
            return
        }

        reader := bufio.NewReader(resp.Body)
        if debug != nil {
            debug.Printf("Starting to read streaming response...")
        }

        for {
            line, err := reader.ReadString('\n')
            if err != nil {
                if err != io.EOF {
                    if debug != nil {
                        debug.Printf("Error reading streaming response: %v", err)
                    }
                    errors <- fmt.Errorf("failed to read response: %w", err)
                } else {
                    if debug != nil {
                        debug.Printf("Reached end of stream")
                    }
                }
                return
            }

            line = strings.TrimSpace(line)
            if !strings.HasPrefix(line, "data: ") {
                continue
            }

            data := strings.TrimPrefix(line, "data: ")
            if data == "[DONE]" {
                if debug != nil {
                    debug.Printf("Received [DONE] from stream")
                }
                return
            }

            if debug != nil {
                debug.Printf("Received chunk: %s", data)
            }

            var streamResp CompletionResponse
            if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
                if debug != nil {
                    debug.Printf("Error parsing chunk: %v", err)
                }
                errors <- fmt.Errorf("failed to parse response: %w", err)
                return
            }

            if streamResp.Usage != nil {
                if debug != nil {
                    debug.Printf("Usage in chunk: %+v", streamResp.Usage)
                }
            }

            // Just pass through whatever content we get
            if len(streamResp.Choices) > 0 {
                responses <- streamResp.Choices[0].Delta.Content
            }
        }
    }()

    return responses, tokens, errors
}