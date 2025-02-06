package main

import (
    "fmt"
    "time"
    "cortex/internal/tokens"
)

func main() {
    examples := []uint64{
        100,    // Small query
        500,    // Medium query
        2048,   // Half context
        4000,   // Almost full context
    }

    for _, tokens := range examples {
        tracker := tokens.NewTracker()
        tracker.AddTokens(tokens)
        time.Sleep(1 * time.Second) // Simulate some processing time
        fmt.Print(tracker.FormatMetrics())
    }
}
