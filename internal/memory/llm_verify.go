package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// verifyMerge checks if a merged element properly represents its source elements
func (p *Processor) verifyMerge(ctx context.Context, original1, original2, merged *Element) (*MergeVerification, error) {
	// Create verification prompt
	verifyPrompt := []Message{
		{
			Role: "system",
			Content: `You are a verification system for merged content. Analyze the original contents and the merged result.
Verify:
1. All critical information is preserved
2. No semantic conflicts exist
3. The merge maintains context
4. The structure is valid
5. The content is logically consistent

Respond with a JSON object containing:
{
    "success": boolean,
    "confidence": float (0-1),
    "issues": [string] (if any),
    "suggestions": [string] (if any)
}`,
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`Verify this merge:
Original 1: %s
Original 2: %s
Merged Result: %s`, original1.Content, original2.Content, merged.Content),
		},
	}

	// Get verification from LLM
	var verificationResponse strings.Builder
	responses, tokens, errors := p.llm.Chat(ctx, verifyPrompt)

	// Collect response chunks
	for content := range responses {
		verificationResponse.WriteString(content)
	}

	// Process tokens if needed
	<-tokens

	// Check for errors
	if err := <-errors; err != nil {
		return nil, fmt.Errorf("LLM verification failed: %w", err)
	}

	// Parse verification result
	var verification MergeVerification
	if err := json.Unmarshal([]byte(verificationResponse.String()), &verification); err != nil {
		return nil, fmt.Errorf("invalid verification response: %w", err)
	}

	return &verification, nil
}