package tokens

import (
    "fmt"
    "sync/atomic"
    "time"
)

// Color codes
const (
    colorCyan = "\033[36m"
    colorReset = "\033[0m"
)

// Tracker keeps track of token usage
type Tracker struct {
    tokens     uint64    // token usage
    maxTokens  uint64    // max tokens for current request
    start      time.Time // request start time
    contextLimit int     // model's context limit
}

// NewTracker creates a new token usage tracker
func NewTracker(contextLimit int) *Tracker {
    return &Tracker{
        start: time.Now(),
        contextLimit: contextLimit,
    }
}

// AddTokens records tokens used
func (t *Tracker) AddTokens(count uint64) {
    atomic.AddUint64(&t.tokens, count)
}

// SetMaxTokens sets the maximum tokens for the current request
func (t *Tracker) SetMaxTokens(max uint64) {
    atomic.StoreUint64(&t.maxTokens, max)
}

// GetTotal returns total tokens used
func (t *Tracker) GetTotal() uint64 {
    return atomic.LoadUint64(&t.tokens)
}

// FormatMetrics returns a terminal-friendly metrics string
func (t *Tracker) FormatMetrics() string {
    tokens := atomic.LoadUint64(&t.tokens)
    percentUsed := float64(tokens) / float64(t.contextLimit) * 100

    contextLimitK := t.contextLimit / 1000
    if t.contextLimit % 1000 != 0 {
        contextLimitK = t.contextLimit / 1000
    }
    return fmt.Sprintf("%sTokens: %d total (%d input + %d output) | %.1f%% of %dk limit | %.1fs%s\n",
        colorCyan,
        tokens,
        tokens/2, // Approximate split between input/output
        tokens/2,
        percentUsed,
        contextLimitK,
        time.Since(t.start).Seconds(),
        colorReset,
    )
}