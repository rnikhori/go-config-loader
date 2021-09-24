package config

import "fmt"

// wrapStage wraps err with a "config: <stage>: " prefix if err is
// non-nil, and returns nil otherwise. It centralizes the repeated error
// wrapping used at each stage of the load pipeline.
func wrapStage(stage string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("config: %s: %w", stage, err)
}
