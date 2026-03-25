package tool

import (
	"context"
)

// Resolver manages multiple tools and provides unified access.
type Resolver interface {
	Resolve(ctx context.Context) ([]Tool, error)
}

// Resolvers is a slice of Resolver that itself implements Resolver.
type Resolvers []Resolver

// Resolve aggregates tools from all resolvers.
func (r Resolvers) Resolve(ctx context.Context) ([]Tool, error) {
	var all []Tool
	for _, resolver := range r {
		tools, err := resolver.Resolve(ctx)
		if err != nil {
			// In a production system, we might log and continue,
			// but here we follow simple aggregate logic.
			continue
		}
		all = append(all, tools...)
	}
	return all, nil
}
