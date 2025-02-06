package memory

import (
    "context"
    "fmt"
    "github.com/google/uuid"
)

func (p *Processor) EnrichContext(ctx context.Context) ([]Element, error) {
    recentElements, err := p.store.GetRecentElements(ctx, 5)
    if err != nil {
        return nil, fmt.Errorf("getting recent elements: %w", err)
    }

    uniqueElements := make(map[uuid.UUID]*Element)
    for _, el := range recentElements {
        uniqueElements[el.ID] = el
    }

    for _, el := range recentElements {
        connected, err := p.store.GetConnectedElements(ctx, el.ID, 2, 0.3)
        if err != nil {
            return nil, fmt.Errorf("getting connected elements: %w", err)
        }

        for _, connEl := range connected {
            uniqueElements[connEl.ID] = connEl
        }
    }

    var enrichedContext []Element
    for _, el := range uniqueElements {
        enrichedContext = append(enrichedContext, *el)
        go func(elementID uuid.UUID) {
            ctx := context.Background()
            p.store.RecordAccess(ctx, elementID, true, 0.0, 0.8)
        }(el.ID)
    }

    return enrichedContext, nil
}