-- Base schema
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS elements (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    content JSONB,
    embedding vector(1536),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS associations (
    id UUID PRIMARY KEY,
    source_id UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,
    target_id UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,
    pattern_type TEXT NOT NULL,
    strength FLOAT NOT NULL CHECK (strength >= 0 AND strength <= 1),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for graph operations
CREATE INDEX IF NOT EXISTS idx_associations_source_id ON associations(source_id);
CREATE INDEX IF NOT EXISTS idx_associations_target_id ON associations(target_id);
CREATE INDEX IF NOT EXISTS idx_associations_source_strength ON associations(source_id, strength);
CREATE INDEX IF NOT EXISTS idx_associations_target_strength ON associations(target_id, strength);
CREATE INDEX IF NOT EXISTS idx_associations_compound ON associations(source_id, target_id, strength);