package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// MergeVerification contains the results of a merge verification
type MergeVerification struct {
	Success     bool      `json:"success"`
	Confidence  float64   `json:"confidence"`
	Issues      []string  `json:"issues,omitempty"`
	Suggestions []string  `json:"suggestions,omitempty"`
}

// MergeElements handles the merging of two similar elements
func (p *Processor) MergeElements(ctx context.Context, element1ID, element2ID uuid.UUID) error {
	// Start transaction
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting merge transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Get both elements
	el1, err := p.store.GetElement(ctx, element1ID)
	if err != nil {
		return fmt.Errorf("getting element1: %w", err)
	}
	el2, err := p.store.GetElement(ctx, element2ID)
	if err != nil {
		return fmt.Errorf("getting element2: %w", err)
	}

	// 2. Create merged content
	mergedContent, err := p.mergeLLMContent(ctx, el1.Content, el2.Content)
	if err != nil {
		return fmt.Errorf("merging content: %w", err)
	}

	// 3. Create new element with merged content
	mergedElement := &Element{
		ID:        uuid.New(),
		Type:      el1.Type, // Elements of same type only
		Content:   mergedContent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 4. Get embedding for merged content
	embedding, err := p.llm.GetEmbedding(ctx, string(mergedContent))
	if err != nil {
		return fmt.Errorf("getting merged embedding: %w", err)
	}
	mergedElement.Embedding = embedding

	// 5. Verify merge
	verification, err := p.verifyMerge(ctx, el1, el2, mergedElement)
	if err != nil {
		return fmt.Errorf("verifying merge: %w", err)
	}

	// If verification failed, store issues and return error
	if !verification.Success {
		issuesJSON, _ := json.Marshal(verification)
		metadataUpdate := fmt.Sprintf(`{"merge_verification_failed": %s}`, issuesJSON)
		if _, err := tx.ExecContext(ctx, `
			UPDATE merge_candidates 
			SET metadata = metadata || $1::jsonb
			WHERE (element1_id = $2 AND element2_id = $3)
			OR (element1_id = $3 AND element2_id = $2)`,
			metadataUpdate, el1.ID, el2.ID); err != nil {
			return fmt.Errorf("updating merge candidate with verification issues: %w", err)
		}
		return fmt.Errorf("merge verification failed: %v", verification.Issues)
	}

	// Add verification results to merged element metadata
	verificationJSON, _ := json.Marshal(verification)
	mergedElement.Content = json.RawMessage(fmt.Sprintf(`{
		"content": %s,
		"merge_verification": %s
	}`, mergedContent, verificationJSON))

	// 6. Create merged element
	if err := p.store.CreateElement(ctx, mergedElement); err != nil {
		return fmt.Errorf("creating merged element: %w", err)
	}

	// 7. Update associations and finalize merge
	if err := p.finalizeMerge(ctx, tx, el1.ID, el2.ID, mergedElement.ID); err != nil {
		return fmt.Errorf("finalizing merge: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing merge transaction: %w", err)
	}

	return nil
}

// finalizeMerge handles the final steps of the merge process
func (p *Processor) finalizeMerge(ctx context.Context, tx *sql.Tx, el1ID, el2ID, mergedID uuid.UUID) error {
	// Update associations
	if err := p.updateAssociations(ctx, tx, el1ID, el2ID, mergedID); err != nil {
		return err
	}

	// Mark original elements as merged
	if err := p.markElementsMerged(ctx, tx, el1ID, el2ID, mergedID); err != nil {
		return err
	}

	// Remove merge candidate
	if err := p.removeMergeCandidate(ctx, tx, el1ID, el2ID); err != nil {
		return err
	}

	return nil
}

// updateAssociations redirects associations to the merged element
func (p *Processor) updateAssociations(ctx context.Context, tx *sql.Tx, el1ID, el2ID, mergedID uuid.UUID) error {
	query := `
		UPDATE associations
		SET source_id = $1
		WHERE source_id IN ($2, $3)
	`
	if _, err := tx.ExecContext(ctx, query, mergedID, el1ID, el2ID); err != nil {
		return err
	}

	query = `
		UPDATE associations
		SET target_id = $1
		WHERE target_id IN ($2, $3)
	`
	if _, err := tx.ExecContext(ctx, query, mergedID, el1ID, el2ID); err != nil {
		return err
	}

	return nil
}

// markElementsMerged marks original elements as merged
func (p *Processor) markElementsMerged(ctx context.Context, tx *sql.Tx, el1ID, el2ID, mergedID uuid.UUID) error {
	metadata := fmt.Sprintf(`{"merged_into": "%s", "merged_at": "%s"}`, 
		mergedID, time.Now().Format(time.RFC3339))

	query := `
		UPDATE elements
		SET content = content || $1::jsonb
		WHERE id IN ($2, $3)
	`
	if _, err := tx.ExecContext(ctx, query, metadata, el1ID, el2ID); err != nil {
		return err
	}

	return nil
}

// removeMergeCandidate removes the processed merge candidate
func (p *Processor) removeMergeCandidate(ctx context.Context, tx *sql.Tx, el1ID, el2ID uuid.UUID) error {
	query := `
		DELETE FROM merge_candidates
		WHERE (element1_id = $1 AND element2_id = $2)
		OR (element1_id = $2 AND element2_id = $1)
	`
	if _, err := tx.ExecContext(ctx, query, el1ID, el2ID); err != nil {
		return err
	}

	return nil
}