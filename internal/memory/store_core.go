package memory

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/google/uuid"
    "github.com/lib/pq"
)

func (s *PostgresStore) CreateElement(ctx context.Context, el *Element) error {
    if el.ID == uuid.Nil {
        el.ID = uuid.New()
    }
    vectorStr := formatVector(el.Embedding)
    _, err := s.db.ExecContext(ctx, `
        INSERT INTO elements (id, type, content, embedding)
        VALUES ($1::uuid, $2, $3, $4::vector)
    `, el.ID, el.Type, el.Content, vectorStr)
    return err
}

func (s *PostgresStore) GetElement(ctx context.Context, id uuid.UUID) (*Element, error) {
    var el Element
    err := s.db.QueryRowContext(ctx, `
        SELECT id, type, content, embedding::float4[], created_at, updated_at
        FROM elements WHERE id = $1::uuid
    `, id).Scan(&el.ID, &el.Type, &el.Content, pq.Array(&el.Embedding),
        &el.CreatedAt, &el.UpdatedAt)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("element not found: %s", id)
    }
    return &el, err
}

func (s *PostgresStore) UpdateElement(ctx context.Context, el *Element) error {
    vectorStr := formatVector(el.Embedding)
    res, err := s.db.ExecContext(ctx, `
        UPDATE elements 
        SET type = $2, content = $3, embedding = $4::vector, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1::uuid
    `, el.ID, el.Type, el.Content, vectorStr)
    if err != nil {
        return fmt.Errorf("failed to update element: %w", err)
    }
    rows, _ := res.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("element not found: %s", el.ID)
    }
    return nil
}

func (s *PostgresStore) DeleteElement(ctx context.Context, id uuid.UUID) error {
    res, err := s.db.ExecContext(ctx, `DELETE FROM elements WHERE id = $1::uuid`, id)
    if err != nil {
        return fmt.Errorf("error deleting element: %w", err)
    }
    rows, _ := res.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("element not found: %s", id)
    }
    return nil
}

func (s *PostgresStore) CreateAssociation(ctx context.Context, assoc *Association) error {
    if assoc.ID == uuid.Nil {
        assoc.ID = uuid.New()
    }
    _, err := s.db.ExecContext(ctx, `
        INSERT INTO associations (id, source_id, target_id, pattern_type, strength, metadata)
        VALUES ($1::uuid, $2::uuid, $3::uuid, $4, $5, $6)
    `, assoc.ID, assoc.SourceID, assoc.TargetID,
        assoc.PatternType, assoc.Strength, assoc.Metadata)
    return err
}