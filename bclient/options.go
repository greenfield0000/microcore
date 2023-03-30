package bclient

import "time"

// clientOptions ...
type clientOptions struct {
	RequestUri     string
	RequestTimeout time.Duration
}
