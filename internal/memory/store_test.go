package memory

import (
    "context"
    "database/sql"
    "encoding/json"
    "testing"
    "time"
    _ "github.com/lib/pq"
    "github.com/google/uuid"
)

func getTestDB() *sql.DB {
    connStr := "host=localhost port=5433 user=cortexai password=cortexai dbname=cortex_dev sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    return db
}

func TestGraphOperations(t *testing.T) {
    ctx := context.Background()
    db := getTestDB()
    store := NewPostgresStore(db)

    content := json.RawMessage(`{"test": "data"}`)
    embedding := make([]float32, 1536)
    metadata := json.RawMessage(`{"test": "metadata"}`)
    
    el1 := &Element{ID: uuid.New(), Type: ElementTypeConcept, Content: content, Embedding: embedding}
    el2 := &Element{ID: uuid.New(), Type: ElementTypeConcept, Content: content, Embedding: embedding}
    el3 := &Element{ID: uuid.New(), Type: ElementTypeConcept, Content: content, Embedding: embedding}

    elements := []*Element{el1, el2, el3}
    for _, el := range elements {
        if err := store.CreateElement(ctx, el); err != nil {
            t.Fatalf("Failed to create element: %v", err)
        }
    }

    assoc1 := &Association{
        SourceID: el1.ID,
        TargetID: el2.ID,
        PatternType: "test",
        Strength: 0.8,
        Metadata: metadata,
    }
    assoc2 := &Association{
        SourceID: el2.ID,
        TargetID: el3.ID,
        PatternType: "test",
        Strength: 0.9,
        Metadata: metadata,
    }
    assoc3 := &Association{
        SourceID: el3.ID,
        TargetID: el1.ID,
        PatternType: "test",
        Strength: 0.7,
        Metadata: metadata,
    }

    associations := []*Association{assoc1, assoc2, assoc3}
    for _, assoc := range associations {
        if err := store.CreateAssociation(ctx, assoc); err != nil {
            t.Fatalf("Failed to create association: %v", err)
        }
    }

    t.Run("FindPaths", func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
        defer cancel()

        paths, err := store.FindPaths(ctx, el1.ID, el3.ID, 3, 0.5)
        if err != nil {
            t.Fatalf("FindPaths failed: %v", err)
        }
        if len(paths) == 0 {
            t.Error("No paths found")
        }
    })

    t.Run("DetectCycles", func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
        defer cancel()

        cycles, err := store.DetectCycles(ctx, el1.ID, 3)
        if err != nil {
            t.Fatalf("DetectCycles failed: %v", err)
        }
        if len(cycles) == 0 {
            t.Error("No cycles detected")
        }
    })

    t.Run("GetConnectedElements", func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
        defer cancel()

        elements, err := store.GetConnectedElements(ctx, el1.ID, 3, 0.5)
        if err != nil {
            t.Fatalf("GetConnectedElements failed: %v", err)
        }
        if len(elements) == 0 {
            t.Error("No connected elements found")
        }
    })
}