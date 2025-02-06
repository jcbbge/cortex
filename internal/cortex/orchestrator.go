package cortex

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "cortex/internal/llm"

    "cortex/internal/tokens"
)

type Orchestrator struct {
    config    *llm.Config
    tracker   *tokens.Tracker
}

func New(config *llm.Config) (*Orchestrator, error) {
    return &Orchestrator{
        config:    config,
        tracker:   tokens.NewTracker(config.GetContextLimit()),
    }, nil
}

func (o *Orchestrator) Process(ctx context.Context, input string) (chan string, *tokens.Tracker, chan error) {
    start := time.Now()
    responses := make(chan string)
    errors := make(chan error)
    preResponses := make(chan string)

    var debug *log.Logger
    if os.Getenv("DEBUG") != "" {
        debug = log.New(os.Stderr, "\n[ORCH] ", log.Ltime|log.Lmicroseconds)
        debug.Printf("Starting request with input length=%d", len(input))
        os.Stderr.Sync()
    }

    go func() {
        defer close(responses)
        defer close(errors)

        // Create client with default model for pre-submission
        modelConfig := *o.config
        modelConfig.Model = "gpt-3.5-turbo"
        client, err := llm.New(&modelConfig)
        if err != nil {
            if debug != nil {
                debug.Printf("Client creation error: %v", err)
            }
            errors <- fmt.Errorf("creating model client: %w", err)
            return
        }

        // Forward pre-submission responses to main response channel
        go func() {
            defer close(preResponses)
            if debug != nil {
                debug.Printf("Starting pre-submission request")
            }
            preResp, preTokens, preErrs := client.Chat(ctx, []llm.Message{
                {Role: "user", Content: input},
            })

            for {
                select {
                case r, ok := <-preResp:
                    if !ok {
                        preResp = nil
                        continue
                    }
                    responses <- r
                case t, ok := <-preTokens:
                    if !ok {
                        preTokens = nil
                        continue
                    }
                    o.tracker.AddTokens(t)
                case e, ok := <-preErrs:
                    if !ok {
                        preErrs = nil
                        continue
                    }
                    errors <- e
                    return
                default:
                    if preResp == nil && preTokens == nil && preErrs == nil {
                        return
                    }
                }
            }
        }()

        // Process with selected model
        if debug != nil {
            debug.Printf("Starting main request")
        }
        resp, tokens, errs := client.Chat(ctx, []llm.Message{
            {Role: "user", Content: input},
        })

        var mainTokenCount uint64
        for {
            select {
            case r, ok := <-resp:
                if !ok {
                    resp = nil
                    continue
                }
                responses <- r
            case t, ok := <-tokens:
                if !ok {
                    tokens = nil
                    continue
                }
                mainTokenCount = t
                if debug != nil {
                    debug.Printf("Received main token count: %d", t)
                }
                o.tracker.AddTokens(t)
                if t > 0 {
                    o.tracker.SetMaxTokens(t) // Only set if we got a valid count
                    if debug != nil {
                        debug.Printf("Set max tokens to %d", t)
                    }
                }
            case e, ok := <-errs:
                if !ok {
                    errs = nil
                    continue
                }
                if debug != nil {
                    debug.Printf("Error from main request: %v", e)
                }
                errors <- e
            default:
                if resp == nil && tokens == nil && errs == nil {
                    if debug != nil {
                        elapsed := time.Since(start)
                        debug.Printf("Request complete - Total=%d Time=%.2fs", 
                            mainTokenCount, elapsed.Seconds())
                    }
                    return
                }
            }
        }
    }()

    return responses, o.tracker, errors
}
