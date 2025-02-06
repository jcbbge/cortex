package memory

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/google/uuid"
)

func (s *PostgresStore) GetRecentElements(ctx context.Context, limit int) ([]*Element, error) {
    return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) FindSimilarElements(ctx context.Context, embedding []float32, limit int, minSimilarity float64) ([]*Element, error) {
    return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetAssociation(ctx context.Context, id uuid.UUID) (*Association, error) {
    return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) UpdateAssociationStrength(ctx context.Context, id uuid.UUID, strengthDelta float64, metadata json.RawMessage) error {
    return fmt.Errorf("not implemented")
}

func (s *PostgresStore) DeleteAssociation(ctx context.Context, id uuid.UUID) error {
    return fmt.Errorf("not implemented")
}

func (s *PostgresStore) RecordAccess(ctx context.Context, elementID uuid.UUID, successful bool, recallTime float64, certaindex float64) error {
    return fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetAccessPattern(ctx context.Context, elementID uuid.UUID) (*AccessPattern, error) {
    return nil, fmt.Errorf("not implemented")
}