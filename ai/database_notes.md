# Database Design Notes
[Previous notes retained...]

## Schema Transition Plan (2025-02-05)

### Current Schema Analysis
```sql
-- Currently we have:
CREATE TABLE elements (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    content JSONB NOT NULL,
    embedding VECTOR(1536),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE associations (
    id UUID PRIMARY KEY,
    source_id UUID NOT NULL REFERENCES elements(id),
    target_id UUID NOT NULL REFERENCES elements(id),
    pattern_type TEXT NOT NULL,
    strength FLOAT NOT NULL DEFAULT 0.0,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE element_paths (
    id UUID PRIMARY KEY,
    element_id UUID NOT NULL REFERENCES elements(id),
    path ltree NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Proposed Schema Changes

1. **Optimize Associations Table**
```sql
-- Add cycle prevention
ALTER TABLE associations 
ADD CONSTRAINT check_no_cycles 
CHECK (NOT has_cycle(source_id, target_id));

-- Add direct path constraint
ALTER TABLE associations
CHECK (source_id != target_id);

-- Add index for recursive traversal
CREATE INDEX idx_associations_traversal 
ON associations(source_id, target_id);
```

2. **Graph Functions**
```sql
-- Cycle detection function
CREATE OR REPLACE FUNCTION has_cycle(
    input_source_id UUID,
    input_target_id UUID
) RETURNS BOOLEAN AS $$
DECLARE
    rec RECORD;
BEGIN
    FOR rec IN
        WITH RECURSIVE traversal AS (
            -- Base case
            SELECT 
                ARRAY[input_source_id] AS path,
                input_target_id AS target_id

            UNION ALL

            -- Recursive case
            SELECT
                traversal.path || a.source_id,
                a.target_id
            FROM traversal
            JOIN associations a ON a.source_id = traversal.target_id
            WHERE NOT a.target_id = ANY(traversal.path)
        )
        SELECT * FROM traversal
    LOOP
        IF rec.target_id = ANY(rec.path) THEN
            RETURN TRUE;
        END IF;
    END LOOP;

    RETURN FALSE;
END;
$$ LANGUAGE plpgsql;

-- Path finding function
CREATE OR REPLACE FUNCTION find_paths(
    start_id UUID,
    end_id UUID
) RETURNS TABLE(
    path UUID[],
    depth INT
) AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE traversal AS (
        -- Base case
        SELECT 
            ARRAY[source_id, target_id] AS path,
            1 AS depth
        FROM associations
        WHERE source_id = start_id

        UNION ALL

        -- Recursive case
        SELECT
            traversal.path || a.target_id,
            traversal.depth + 1
        FROM traversal
        JOIN associations a ON a.source_id = traversal.path[array_length(traversal.path, 1)]
        WHERE NOT a.target_id = ANY(traversal.path)
    )
    SELECT path, depth 
    FROM traversal
    WHERE path[array_length(path, 1)] = end_id;
END;
$$ LANGUAGE plpgsql;
```

3. **Migration Strategy**
- Phase 1: Add new functions and constraints
- Phase 2: Migrate data from element_paths if needed
- Phase 3: Remove element_paths table (optional)

### Performance Considerations
1. **Indexing Strategy**
   - Primary key on associations
   - Composite index for traversal
   - Vector index maintained on elements
   - Pattern type index for filtering

2. **Query Optimization**
   - Early stopping in recursive queries
   - Path array for cycle prevention
   - Depth limits for safety
   - Transaction boundaries for consistency

3. **Memory Considerations**
   - Path arrays grow with depth
   - Consider depth limits based on use case
   - Monitor recursive query performance
   - Cache frequent paths if needed

### Next Steps
1. Implement cycle detection function
2. Add constraints to associations table
3. Test path finding implementation
4. Benchmark recursive queries
5. Monitor memory usage patterns

### Future Optimizations
1. Consider materialized paths for frequent queries
2. Implement path caching if needed
3. Add statistics gathering
4. Tune recursive query parameters

### Questions to Address
1. Maximum practical depth for recursion?
2. Caching strategy for frequent paths?
3. Monitoring approach for query performance?
4. Backup strategy during migration?