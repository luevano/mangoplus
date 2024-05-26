package creators

import (
	"fmt"
)

const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"

type Options struct {
	UserAgent string
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
