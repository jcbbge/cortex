-- Indexes for graph operations

-- Association traversal
CREATE INDEX IF NOT EXISTS idx_associations_source_id ON associations(source_id);
CREATE INDEX IF NOT EXISTS idx_associations_target_id ON associations(target_id);

-- Combined index for association strength filtering
CREATE INDEX IF NOT EXISTS idx_associations_source_strength ON associations(source_id, strength);
CREATE INDEX IF NOT EXISTS idx_associations_target_strength ON associations(target_id, strength);

-- Speed up recursive CTEs
CREATE INDEX IF NOT EXISTS idx_associations_compound ON associations(source_id, target_id, strength);

ANALYZE associations;