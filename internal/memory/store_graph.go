package memory

import (
    "context"
    "fmt"
    "github.com/google/uuid"
    "github.com/lib/pq"
)

func (s *PostgresStore) GetConnectedElements(ctx context.Context, elementID uuid.UUID, maxDepth int, minStrength float64) ([]*Element, error) {
    rows, err := s.db.QueryContext(ctx, `
    WITH RECURSIVE connected_ids AS (
        -- Base case: direct connections
        SELECT DISTINCT
            CASE WHEN a.source_id = $1::uuid THEN a.target_id ELSE a.source_id END as id,
            ARRAY[CASE WHEN a.source_id = $1::uuid THEN a.target_id ELSE a.source_id END] as path,
            1 as depth
        FROM associations a
        WHERE (a.source_id = $1::uuid OR a.target_id = $1::uuid)
            AND a.strength >= $3

        UNION

        -- Recursive case
        SELECT DISTINCT
            CASE WHEN a.source_id = c.id THEN a.target_id ELSE a.source_id END,
            c.path || CASE WHEN a.source_id = c.id THEN a.target_id ELSE a.source_id END,
            c.depth + 1
        FROM connected_ids c
        JOIN associations a ON (a.source_id = c.id OR a.target_id = c.id)
        WHERE c.depth < $2
            AND a.strength >= $3
            AND CASE WHEN a.source_id = c.id THEN a.target_id ELSE a.source_id END != $1::uuid
            AND NOT CASE WHEN a.source_id = c.id THEN a.target_id ELSE a.source_id END = ANY(c.path)
    )
    SELECT DISTINCT ON (e.id)
        e.id, e.type, e.content, e.embedding::float4[], e.created_at, e.updated_at
    FROM connected_ids c
    JOIN elements e ON e.id = c.id
    ORDER BY e.id, c.depth;
    `, elementID, maxDepth, minStrength)
    if err != nil {
        return nil, fmt.Errorf("error getting connected elements: %w", err)
    }
    defer rows.Close()

    var elements []*Element
    for rows.Next() {
        var el Element
        err := rows.Scan(
            &el.ID,
            &el.Type,
            &el.Content,
            pq.Array(&el.Embedding),
            &el.CreatedAt,
            &el.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning element: %w", err)
        }
        elements = append(elements, &el)
    }
    return elements, nil
}