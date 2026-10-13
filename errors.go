package mangoplus

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Errors to check with errors.Is, the error types below match them when it applies.
//
// ErrRateLimited (a 429, or the API's 10522) has no retry delay on purpose: 10522 reportedly locks
// the device out for 10+ minutes, retrying within seconds can't work and may keep it locked, so
// callers should stop and tell the user instead of waiting.
var (
	ErrNotFound             = errors.New("not found")
	ErrRateLimited          = errors.New("rate limited")
	ErrMaintenance          = errors.New("under maintenance")
	ErrSubscriptionRequired = errors.New("subscription required")
	ErrRegionBlocked        = errors.New("not available in this region")
)

// HTTPError: The API answered with a non-200 status code.
type HTTPError struct {
	Method     string
	Path       string
	StatusCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("mangoplus: %s %s: %d %s", e.Method, e.Path, e.StatusCode, strings.ToLower(http.StatusText(e.StatusCode)))
}

func (e *HTTPError) Is(target error) bool {
	switch e.StatusCode {
	case http.StatusNotFound:
		return target == ErrNotFound
	case http.StatusTooManyRequests:
		return target == ErrRateLimited
	case http.StatusServiceUnavailable:
		return target == ErrMaintenance
	}
	return false
}

// ServerError: The API answered with an error result, e.g. "Invalid user access(11301)".
type ServerError struct {
	Method string
	Path   string
	Action ErrorAction
	// Code: Number at the end of the message, 0 if there's none.
	Code    int
	Message string
}

func (e *ServerError) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("mangoplus: %s %s: %s (%d)", e.Method, e.Path, e.Message, e.Code)
	}
	return fmt.Sprintf("mangoplus: %s %s: %s", e.Method, e.Path, e.Message)
}

func (e *ServerError) Is(target error) bool {
	switch e.Action {
	case ErrorAction_MAINTENANCE:
		return target == ErrMaintenance
	case ErrorAction_GEOIP_BLOCKING:
		return target == ErrRegionBlocked
	}
	switch e.Code {
	case 11301:
		return target == ErrSubscriptionRequired
	case 10522:
		return target == ErrRateLimited
	case 11103:
		return target == ErrNotFound
	case 11101:
		// a generic "failed to load the page", on a title's details it's an unknown title
		return target == ErrNotFound && strings.Contains(e.Path, MangaPath)
	}
	return false
}

// the code comes at the end of the popup body, e.g. "Invalid user access(11301)"
var serverErrorCode = regexp.MustCompile(`\s*\((\d+)\)$`)

func newServerError(method, path string, res *ErrorResult) *ServerError {
	popup := res.EnglishPopup
	if popup == nil {
		popup = res.SpanishPopup
	}
	if popup == nil && len(res.Popups) != 0 {
		popup = res.Popups[0]
	}

	var message string
	switch {
	case popup != nil && popup.Body != "":
		message = popup.Body
	case popup != nil:
		message = popup.Subject
	case res.DebugInfo != "":
		message = res.DebugInfo
	default:
		message = "error action " + res.Action.String()
	}

	var code int
	if match := serverErrorCode.FindStringSubmatch(message); match != nil {
		code, _ = strconv.Atoi(match[1])
		message = strings.TrimSuffix(message, match[0])
	}
	return &ServerError{Method: method, Path: path, Action: res.Action, Code: code, Message: message}
}
