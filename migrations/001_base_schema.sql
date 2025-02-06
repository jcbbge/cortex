-- Base schema for Cortex v0.5 LLM-native associative memory

-- Elements table - stores base units of knowledge
CREATE TABLE IF NOT EXISTS elements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL CHECK (type IN ('code', 'concept', 'context')),
    content JSONB NOT NULL,
    embedding vector(1536),  -- NULL until computed
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Associations table - tracks LLM-discovered patterns
CREATE TABLE IF NOT EXISTS associations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,
    target_id UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,
    pattern_type TEXT NOT NULL,
    strength FLOAT NOT NULL DEFAULT 0.0 CHECK (strength >= 0.0 AND strength <= 1.0),
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Access patterns table - tracks usage and performance
CREATE TABLE IF NOT EXISTS access_patterns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    element_id UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,
    access_count INTEGER NOT NULL DEFAULT 0,
    successful_recalls INTEGER NOT NULL DEFAULT 0,
    avg_recall_time FLOAT,
    last_accessed_at TIMESTAMPTZ,
    certaindex FLOAT CHECK (certaindex >= 0.0 AND certaindex <= 1.0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Element paths table - DAG structure for relationships
CREATE TABLE IF NOT EXISTS element_paths (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    element_id UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,
    path ltree NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS elements_type_idx ON elements(type);
CREATE INDEX IF NOT EXISTS elements_embedding_idx ON elements USING ivfflat (embedding vector_cosine_ops);

CREATE INDEX IF NOT EXISTS associations_source_idx ON associations(source_id);
CREATE INDEX IF NOT EXISTS associations_target_idx ON associations(target_id);
CREATE INDEX IF NOT EXISTS associations_pattern_type_idx ON associations(pattern_type);

CREATE INDEX IF NOT EXISTS access_patterns_element_idx ON access_patterns(element_id);
CREATE INDEX IF NOT EXISTS element_paths_path_idx ON element_paths USING GIST (path);

-- Helper function to update timestamps
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Update triggers
CREATE TRIGGER update_elements_timestamp
    BEFORE UPDATE ON elements
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_associations_timestamp
    BEFORE UPDATE ON associations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_access_patterns_timestamp
    BEFORE UPDATE ON access_patterns
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();