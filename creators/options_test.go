package creators

import (
	"net/http"
	"testing"
)

func TestHTTPClientOption(t *testing.T) {
	c, err := NewCreatorsClient(DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if c.client.Timeout == 0 {
		t.Error("the default client should have a timeout")
	}

	given := &http.Client{}
	o := DefaultOptions()
	o.HTTPClient = given
	if c, err = NewCreatorsClient(o); err != nil {
		t.Fatal(err)
	}
	if c.client != given {
		t.Error("the given client should be used")
	}
}
