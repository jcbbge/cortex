package memory

import (
    "context"
    "encoding/json"
    "time"
    "github.com/google/uuid"
)

type ElementType string

const (
    ElementTypeCode    ElementType = "code"
    ElementTypeConcept ElementType = "concept"
    ElementTypeContext ElementType = "context"
)

type Element struct {
    ID        uuid.UUID       `json:"id"`
    Type      ElementType     `json:"type"`
    Content   json.RawMessage `json:"content"`
    Embedding []float32       `json:"embedding,omitempty"`
    CreatedAt time.Time       `json:"created_at"`
    UpdatedAt time.Time       `json:"updated_at"`
}

type Association struct {
    ID          uuid.UUID       `json:"id"`
    SourceID    uuid.UUID       `json:"source_id"`
    TargetID    uuid.UUID       `json:"target_id"`
    PatternType string          `json:"pattern_type"`
    Strength    float64         `json:"strength"`
    Metadata    json.RawMessage `json:"metadata,omitempty"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
}

type PathNode struct {
    ElementID uuid.UUID
    Depth     int
    Path      []uuid.UUID
    Strength  float64
}

type AccessPattern struct {
    ID              uuid.UUID  `json:"id"`
    ElementID       uuid.UUID  `json:"element_id"`
    AccessCount     int        `json:"access_count"`
    SuccessfulCount int        `json:"successful_recalls"`
    AvgRecallTime   float64    `json:"avg_recall_time"`
    LastAccessedAt  time.Time  `json:"last_accessed_at"`
    Certaindex      float64    `json:"certaindex"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
}

type MemoryStore interface {
    CreateElement(ctx context.Context, el *Element) error
    GetElement(ctx context.Context, id uuid.UUID) (*Element, error)
    UpdateElement(ctx context.Context, el *Element) error
    DeleteElement(ctx context.Context, id uuid.UUID) error
    GetRecentElements(ctx context.Context, limit int) ([]*Element, error)
    FindSimilarElements(ctx context.Context, embedding []float32, limit int, minSimilarity float64) ([]*Element, error)
    
    CreateAssociation(ctx context.Context, assoc *Association) error
    GetAssociation(ctx context.Context, id uuid.UUID) (*Association, error)
    UpdateAssociationStrength(ctx context.Context, id uuid.UUID, strengthDelta float64, metadata json.RawMessage) error
    DeleteAssociation(ctx context.Context, id uuid.UUID) error

    RecordAccess(ctx context.Context, elementID uuid.UUID, successful bool, recallTime float64, certaindex float64) error
    GetAccessPattern(ctx context.Context, elementID uuid.UUID) (*AccessPattern, error)

    FindPaths(ctx context.Context, sourceID, targetID uuid.UUID, maxDepth int, minStrength float64) ([][]PathNode, error)
    GetConnectedElements(ctx context.Context, elementID uuid.UUID, maxDepth int, minStrength float64) ([]*Element, error)
    DetectCycles(ctx context.Context, elementID uuid.UUID, maxDepth int) ([][]uuid.UUID, error)
}