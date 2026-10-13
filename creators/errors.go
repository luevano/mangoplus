package creators

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ErrNotFound: Check with errors.Is.
var ErrNotFound = errors.New("not found")

// HTTPError: The API answered with an unexpected status code.
type HTTPError struct {
	Method     string
	Path       string
	StatusCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("mangapluscreators: %s %s: %d %s", e.Method, e.Path, e.StatusCode, strings.ToLower(http.StatusText(e.StatusCode)))
}

func (e *HTTPError) Is(target error) bool {
	return e.StatusCode == http.StatusNotFound && target == ErrNotFound
}

// ErrBlocked: AWS WAF answered instead of the API, with a challenge only a browser can solve or
// a block. Check with errors.Is.
var ErrBlocked = errors.New("blocked by AWS WAF")
