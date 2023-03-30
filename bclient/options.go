package bclient

import "time"

// ClientOptions ...
type ClientOptions struct {
	RequestUri     string
	RequestTimeout time.Duration
}
