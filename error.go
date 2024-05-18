package mangoplus

import (
	"fmt"
	"strings"
)

type Popup struct {
	Subject      string `json:"subject"`
	Body         string `json:"body"`
	CancelButton struct {
		Text string `json:"text"`
	} `json:"cancelButton"`
	Language *string `json:"language"`
}


// TODO: handle one specific language instead of the first in the list
//
// GetErrors: Get the errors for this particular request.
func (error *ErrorResponse) GetErrors() string {
	var errors strings.Builder
	var errorPopup Popup
	if error.EnglishPopup != nil {
		errorPopup = *error.EnglishPopup
	} else if error.SpanishPopup != nil {
		errorPopup = *error.SpanishPopup
	} else {
		for _, err := range *error.Popups {
			errorPopup = err
			break
		}
	}
	errors.WriteString(fmt.Sprintf("%s: %s", errorPopup.Subject, errorPopup.Body))
	return errors.String()
}
