package mangoplus

import "fmt"

// message: Get the error message, preferring the english popup.
func (e *ErrorResult) message() string {
	popup := e.EnglishPopup
	if popup == nil {
		popup = e.SpanishPopup
	}
	if popup == nil && len(e.Popups) != 0 {
		popup = e.Popups[0]
	}
	if popup != nil {
		return fmt.Sprintf("%s: %s", popup.Subject, popup.Body)
	}
	if e.DebugInfo != "" {
		return e.DebugInfo
	}
	return fmt.Sprintf("action %s", e.Action)
}
