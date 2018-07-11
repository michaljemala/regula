package store

import (
	"context"
	"errors"

	"github.com/heetch/regula/rule"
)

// Errors
var (
	ErrNotFound = errors.New("not found")
)

// Store manages the storage of rulesets.
type Store interface {
	// List returns all the rulesets entries under the given prefix.
	List(ctx context.Context, prefix string) ([]RulesetEntry, error)
	// One returns the ruleset entry which corresponds to the given path.
	One(ctx context.Context, path string) (*RulesetEntry, error)
	Watch(ctx context.Context, prefix string) ([]Event, error)
}

// RulesetEntry holds a ruleset and its metadata.
type RulesetEntry struct {
	Path    string
	Ruleset *rule.Ruleset
}

// List of possible events executed against a ruleset.
const (
	PutEvent    = "PUT"
	DeleteEvent = "DELETE"
)

// Event describes an event that occured on a ruleset.
type Event struct {
	Type    string
	Path    string
	Ruleset *rule.Ruleset
}
