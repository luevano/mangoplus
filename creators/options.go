package creators

import (
	"fmt"
	"net/http"
	"time"
)

const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"

// defaultTimeout: Timeout of the client used when Options.HTTPClient is nil.
const defaultTimeout = time.Minute

type Options struct {
	UserAgent string

	// HTTPClient: Client used for every request.
	//
	// Optional, a client with a 1 minute timeout is used when nil.
	HTTPClient *http.Client
}

func (o Options) validate() error {
	if o.UserAgent == "" {
		return fmt.Errorf("UserAgent is empty")
	}
	return nil
}

func DefaultOptions() Options {
	return Options{UserAgent: defaultUserAgent}
}
