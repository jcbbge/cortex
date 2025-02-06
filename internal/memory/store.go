package memory

import (
    "context"
    "database/sql"
    "fmt"
    "strings"

    "github.com/google/uuid"
    "github.com/lib/pq"
)

type PostgresStore struct {
    db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
    return &PostgresStore{db: db}
}

func formatVector(v []float32) string {
    nums := make([]string, len(v))
    for i, x := range v {
        nums[i] = fmt.Sprintf("%f", x)
    }
    return "[" + strings.Join(nums, ",") + "]"
}

func (s *PostgresStore) FindPaths(ctx context.Context, sourceID, targetID uuid.UUID, maxDepth int, minStrength float64) ([][]PathNode, error) {
	rows, err := s.db.QueryContext(ctx, `
		WITH RECURSIVE paths AS (
			-- Base case: direct connections
			SELECT 
				source_id,
				target_id,
				ARRAY[source_id, target_id] as path,
				1 as depth
			FROM associations
			WHERE strength >= $3
			UNION ALL
			-- Recursive case
			SELECT
				p.source_id,
				a.target_id,
				path || a.target_id,
				p.depth + 1
			FROM paths p
			JOIN associations a ON p.target_id = a.source_id
			WHERE a.strength >= $3
			AND p.depth < $4
			AND NOT a.target_id = ANY(p.path)
		)
		SELECT DISTINCT path
		FROM paths
		WHERE source_id = $1 AND target_id = $2
	`, sourceID, targetID, minStrength, maxDepth)
	if err != nil {
		return nil, fmt.Errorf("error finding paths: %v", err)
	}
	defer rows.Close()

	var paths [][]PathNode
	for rows.Next() {
		var path []uuid.UUID
		if err := rows.Scan(pq.Array(&path)); err != nil {
			return nil, fmt.Errorf("error scanning path: %v", err)
		}
		
		// Convert UUIDs to PathNodes
		pathNodes := make([]PathNode, len(path))
		for i, id := range path {
			pathNodes[i] = PathNode{ElementID: id}
		}
		paths = append(paths, pathNodes)
	}

	return paths, nil
}

func (s *PostgresStore) DetectCycles(ctx context.Context, elementID uuid.UUID, maxDepth int) ([][]uuid.UUID, error) {
	rows, err := s.db.QueryContext(ctx, `
		WITH RECURSIVE cycles AS (
			-- Base case: start from the element
			SELECT 
				source_id,
				target_id,
				ARRAY[source_id] as path,
				0 as depth
			FROM associations
			WHERE source_id = $1
			UNION ALL
			-- Recursive case
			SELECT
				c.source_id,
				a.target_id,
				path || a.target_id,
				c.depth + 1
			FROM cycles c
			JOIN associations a ON c.target_id = a.source_id
			WHERE c.depth < $2
			AND (a.target_id = $1 OR NOT a.target_id = ANY(c.path))
		)
		SELECT DISTINCT path || target_id
		FROM cycles
		WHERE target_id = $1 AND array_length(path, 1) > 1
	`, elementID, maxDepth)
	if err != nil {
		return nil, fmt.Errorf("error detecting cycles: %v", err)
	}
	defer rows.Close()

	var cycles [][]uuid.UUID
	for rows.Next() {
		var cycle []uuid.UUID
		if err := rows.Scan(pq.Array(&cycle)); err != nil {
			return nil, fmt.Errorf("error scanning cycle: %v", err)
		}
		cycles = append(cycles, cycle)
	}

	return cycles, nil
}