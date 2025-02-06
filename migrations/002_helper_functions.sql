-- Helper functions for Cortex v0.5 memory system

-- Create or update an element with its embedding
CREATE OR REPLACE FUNCTION upsert_element(
    p_type TEXT,
    p_content JSONB,
    p_embedding vector(1536) DEFAULT NULL
) RETURNS UUID AS $$
DECLARE
    v_id UUID;
BEGIN
    INSERT INTO elements (type, content, embedding)
    VALUES (p_type, p_content, p_embedding)
    RETURNING id INTO v_id;

    -- Initialize access pattern entry
    INSERT INTO access_patterns (element_id)
    VALUES (v_id);

    RETURN v_id;
END;
$$ LANGUAGE plpgsql;

-- Create an association between elements
CREATE OR REPLACE FUNCTION create_association(
    p_source_id UUID,
    p_target_id UUID,
    p_pattern_type TEXT,
    p_initial_strength FLOAT DEFAULT 0.5,
    p_metadata JSONB DEFAULT '{}'::JSONB
) RETURNS UUID AS $$
DECLARE
    v_id UUID;
BEGIN
    INSERT INTO associations (source_id, target_id, pattern_type, strength, metadata)
    VALUES (p_source_id, p_target_id, p_pattern_type, p_initial_strength, p_metadata)
    RETURNING id INTO v_id;

    RETURN v_id;
END;
$$ LANGUAGE plpgsql;

-- Record an access to an element
CREATE OR REPLACE FUNCTION record_element_access(
    p_element_id UUID,
    p_successful BOOLEAN DEFAULT true,
    p_recall_time FLOAT DEFAULT NULL,
    p_certaindex FLOAT DEFAULT NULL
) RETURNS VOID AS $$
BEGIN
    UPDATE access_patterns
    SET 
        access_count = access_count + 1,
        successful_recalls = successful_recalls + CASE WHEN p_successful THEN 1 ELSE 0 END,
        avg_recall_time = CASE 
            WHEN avg_recall_time IS NULL THEN p_recall_time
            ELSE (avg_recall_time * access_count + p_recall_time) / (access_count + 1)
        END,
        last_accessed_at = NOW(),
        certaindex = COALESCE(p_certaindex, certaindex)
    WHERE element_id = p_element_id;
END;
$$ LANGUAGE plpgsql;

-- Update association strength
CREATE OR REPLACE FUNCTION update_association_strength(
    p_association_id UUID,
    p_strength_delta FLOAT,
    p_metadata_update JSONB DEFAULT NULL
) RETURNS FLOAT AS $$
DECLARE
    v_new_strength FLOAT;
BEGIN
    UPDATE associations
    SET 
        strength = GREATEST(0.0, LEAST(1.0, strength + p_strength_delta)),
        metadata = CASE 
            WHEN p_metadata_update IS NOT NULL 
            THEN metadata || p_metadata_update 
            ELSE metadata 
        END
    WHERE id = p_association_id
    RETURNING strength INTO v_new_strength;

    RETURN v_new_strength;
END;
$$ LANGUAGE plpgsql;

-- Find similar elements by embedding
CREATE OR REPLACE FUNCTION find_similar_elements(
    p_embedding vector(1536),
    p_limit INTEGER DEFAULT 5,
    p_min_similarity FLOAT DEFAULT 0.7
) RETURNS TABLE (
    id UUID,
    type TEXT,
    content JSONB,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        e.id,
        e.type,
        e.content,
        (e.embedding <=> p_embedding) as similarity
    FROM elements e
    WHERE e.embedding IS NOT NULL
    AND (e.embedding <=> p_embedding) >= p_min_similarity
    ORDER BY similarity DESC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;

-- Get element context through associations
CREATE OR REPLACE FUNCTION get_element_context(
    p_element_id UUID,
    p_max_depth INTEGER DEFAULT 2,
    p_min_strength FLOAT DEFAULT 0.3
) RETURNS TABLE (
    source_id UUID,
    target_id UUID,
    pattern_type TEXT,
    strength FLOAT,
    depth INTEGER
) AS $$
WITH RECURSIVE context_tree AS (
    -- Base case
    SELECT 
        a.source_id,
        a.target_id,
        a.pattern_type,
        a.strength,
        1 as depth
    FROM associations a
    WHERE (a.source_id = p_element_id OR a.target_id = p_element_id)
    AND a.strength >= p_min_strength

    UNION

    -- Recursive case
    SELECT 
        a.source_id,
        a.target_id,
        a.pattern_type,
        a.strength,
        ct.depth + 1
    FROM associations a
    INNER JOIN context_tree ct ON (a.source_id = ct.target_id OR a.target_id = ct.source_id)
    WHERE a.strength >= p_min_strength
    AND ct.depth < p_max_depth
)
SELECT DISTINCT * FROM context_tree;
$$ LANGUAGE sql;