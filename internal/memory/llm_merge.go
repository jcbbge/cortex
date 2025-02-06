package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// mergeLLMContent uses LLM to intelligently merge content
func (p *Processor) mergeLLMContent(ctx context.Context, content1, content2 json.RawMessage) (json.RawMessage, error) {
	// Create merge prompt
	mergePrompt := []Message{
		{
			Role: "system",
			Content: `You are an AI tasked with merging similar content while preserving important information and removing redundancy.
Rules:
1. Combine overlapping information
2. Preserve unique details from both sources
3. Maintain the original structure where possible
4. Format the output as valid JSON
5. Ensure semantic consistency`,
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`Merge these two pieces of content:
Content 1: %s
Content 2: %s
Respond with only the merged content in the same JSON structure as the inputs.`, content1, content2),
		},
	}

	// Get merged content from LLM
	var mergedContent strings.Builder
	responses, tokens, errors := p.llm.Chat(ctx, mergePrompt)

	// Collect all response chunks
	for content := range responses {
		mergedContent.WriteString(content)
	}

	// Process tokens if needed
	<-tokens

	// Check for errors
	if err := <-errors; err != nil {
		return nil, fmt.Errorf("LLM merge failed: %w", err)
	}

	// Validate JSON structure
	var jsonCheck interface{}
	if err := json.Unmarshal([]byte(mergedContent.String()), &jsonCheck); err != nil {
		return nil, fmt.Errorf("merged content is not valid JSON: %w", err)
	}

	return json.RawMessage(mergedContent.String()), nil
}