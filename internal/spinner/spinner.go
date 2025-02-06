package spinner

import (
    "fmt"
    "time"
)

type Spinner struct {
    stopChan chan struct{}
}

func New() *Spinner {
    return &Spinner{
        stopChan: make(chan struct{}),
    }
}

func (s *Spinner) Start() {
    frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
    go func() {
        for i := 0; ; i++ {
            select {
            case <-s.stopChan:
                fmt.Print("\r") // Clear spinner
                return
            default:
                fmt.Printf("\r%s ", frames[i%len(frames)]) // Add space after spinner
                time.Sleep(100 * time.Millisecond)
            }
        }
    }()
}

func (s *Spinner) Stop() {
    select {
    case s.stopChan <- struct{}{}:
    default:
        // Channel is already closed or no receiver
    }
    fmt.Print("\r") // Clear spinner
}