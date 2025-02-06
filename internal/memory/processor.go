package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Processor handles memory operations and optimization
type Processor struct {
	db    *sql.DB
	store MemoryStore
	llm   LLMClient
}

// NewProcessor creates a new memory processor
func NewProcessor(db *sql.DB, store MemoryStore, llm LLMClient) *Processor {
	return &Processor{
		db:    db,
		store: store,
		llm:   llm,
	}
}

// ProcessInput processes user input and stores it in memory
func (p *Processor) ProcessInput(ctx context.Context, input string) error {
	// 1. Create an element for the input
	inputEl := &Element{
		ID:        uuid.New(),
		Type:      ElementTypeContext,
		Content:   json.RawMessage(input),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 2. Get embedding for the input
	embedding, err := p.llm.GetEmbedding(ctx, input)
	if err != nil {
		return err
	}
	inputEl.Embedding = embedding

	// 3. Store the element
	if err := p.store.CreateElement(ctx, inputEl); err != nil {
		return err
	}

	// 4. Find similar elements
	similar, err := p.store.FindSimilarElements(ctx, embedding, 5, 0.7)
	if err != nil {
		return err
	}

	// 5. Create associations with similar elements
	for _, el := range similar {
		assoc := &Association{
			ID:          uuid.New(),
			SourceID:    inputEl.ID,
			TargetID:    el.ID,
			PatternType: "embedding_similarity",
			Strength:    0.5, // Initial strength
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := p.store.CreateAssociation(ctx, assoc); err != nil {
			return err
		}
	}

	return nil
}