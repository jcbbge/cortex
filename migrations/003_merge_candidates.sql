-- Migration for merge_candidates table
-- Create merge_candidates table for tracking potential element merges
CREATE TABLE IF NOT EXISTS merge_candidates (
    element1_id UUID NOT NULL REFERENCES elements(id),
    element2_id UUID NOT NULL REFERENCES elements(id),
    similarity FLOAT NOT NULL CHECK (similarity > 0 AND similarity <= 1),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (element1_id, element2_id)
);

-- Add index for querying candidates by similarity threshold
CREATE INDEX merge_candidates_similarity_idx ON merge_candidates(similarity DESC);

-- Add index for checking existing candidates quickly
CREATE INDEX merge_candidates_elements_idx ON merge_candidates(element2_id, element1_id);

-- Function to cleanup old merge candidates
CREATE OR REPLACE FUNCTION cleanup_old_merge_candidates() RETURNS void AS $$
BEGIN
    -- Remove candidates older than 7 days that haven't been merged
    DELETE FROM merge_candidates
    WHERE created_at < NOW() - INTERVAL '7 days';
END;
$$ LANGUAGE plpgsql;

-- Comments for maintainability
COMMENT ON TABLE merge_candidates IS 'Tracks potential element pairs for merging based on similarity';
COMMENT ON COLUMN merge_candidates.element1_id IS 'ID of the first element in the merge pair';
COMMENT ON COLUMN merge_candidates.element2_id IS 'ID of the second element in the merge pair';
COMMENT ON COLUMN merge_candidates.similarity IS 'Cosine similarity between the elements (0-1)';
COMMENT ON COLUMN merge_candidates.created_at IS 'When this merge candidate was identified';