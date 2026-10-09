package mangoplus

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
)

// newErrorTestClient: Client pointing to a local server that answers every request with the given
// status code, and the given response body if any.
func newErrorTestClient(t *testing.T, status int, res *Response) *PlusClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if res != nil {
			body, _ := proto.Marshal(res)
			w.Write(body)
		}
	}))
	t.Cleanup(server.Close)

	o := DefaultOptions()
	o.Secret = testSecret // skip registering
	c, err := NewPlusClient(o)
	if err != nil {
		t.Fatal(err)
	}
	c.baseAPI = server.URL
	return c
}

func popupError(action ErrorAction, body string) *Response {
	return &Response{Error: &ErrorResult{Action: action, EnglishPopup: &Popup{Subject: "subject", Body: body}}}
}

func TestHTTPError(t *testing.T) {
	c := newErrorTestClient(t, http.StatusForbidden, nil)

	_, err := c.Manga.Get("100310")
	want := `get manga "100310": mangoplus: GET /title_detailV3: 403 forbidden`
	if err == nil || err.Error() != want {
		t.Fatalf("wanted %q, got %v", want, err)
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusForbidden {
		t.Errorf("wanted an HTTPError with status 403, got %#v", httpErr)
	}
}

func TestServerError(t *testing.T) {
	tests := []struct {
		res      *Response
		sentinel error
		message  string
	}{
		{popupError(ErrorAction_DEFAULT, "Invalid user access(11301)"), ErrSubscriptionRequired,
			`get chapter pages "1": mangoplus: GET /manga_viewer_v3: Invalid user access (11301)`},
		{popupError(ErrorAction_DEFAULT, "Invalid Parameter(10522)"), ErrRateLimited,
			`get chapter pages "1": mangoplus: GET /manga_viewer_v3: Invalid Parameter (10522)`},
		{popupError(ErrorAction_MAINTENANCE, "Under maintenance"), ErrMaintenance,
			`get chapter pages "1": mangoplus: GET /manga_viewer_v3: Under maintenance`},
		{popupError(ErrorAction_GEOIP_BLOCKING, "Not available"), ErrRegionBlocked,
			`get chapter pages "1": mangoplus: GET /manga_viewer_v3: Not available`},
	}
	for _, tt := range tests {
		c := newErrorTestClient(t, http.StatusOK, tt.res)

		_, err := c.Page.Get("1", false, ImageQualityHigh)
		if err == nil || err.Error() != tt.message {
			t.Errorf("wanted %q, got %v", tt.message, err)
		}
		if !errors.Is(err, tt.sentinel) {
			t.Errorf("%v should match %v", err, tt.sentinel)
		}
	}
}

func TestServerErrorFromTestdata(t *testing.T) {
	body, err := os.ReadFile("testdata/error.bin")
	if err != nil {
		t.Fatal(err)
	}
	var res Response
	if err := proto.Unmarshal(body, &res); err != nil {
		t.Fatal(err)
	}
	c := newErrorTestClient(t, http.StatusOK, &res)

	_, err = c.Page.Get("1", false, ImageQualityHigh)
	var serverErr *ServerError
	if !errors.As(err, &serverErr) || serverErr.Code != 11103 || !errors.Is(err, ErrNotFound) {
		t.Errorf("wanted a not found ServerError with code 11103, got %v", err)
	}
}

func TestInvalidOptions(t *testing.T) {
	o := DefaultOptions()
	o.UserAgent = ""
	if _, err := NewPlusClient(o); err == nil || err.Error() != "mangoplus: invalid options: UserAgent is empty" {
		t.Errorf("wanted an invalid options error, got %v", err)
	}
}

// newClosedTestClient: Client pointing to a closed local server, so requests fail at the network level.
func newClosedTestClient(t *testing.T, options Options) *PlusClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()
	c, err := NewPlusClient(options)
	if err != nil {
		t.Fatal(err)
	}
	c.baseAPI = server.URL
	return c
}

func TestSecretsRedactedFromErrors(t *testing.T) {
	o := DefaultOptions()
	o.Secret = "paid-secret"
	o.FreeSecret = "free-secret"
	_, err := newClosedTestClient(t, o).Manga.Get("1")
	if err == nil || strings.Contains(err.Error(), "paid-secret") || !strings.Contains(err.Error(), "secret=REDACTED") {
		t.Errorf("secret should be redacted, got %v", err)
	}

	// registering sends the device token and security key
	o = DefaultOptions()
	deviceToken := md5Hex(o.AndroidID)
	_, err = newClosedTestClient(t, o).Manga.Get("1")
	if err == nil || strings.Contains(err.Error(), deviceToken) || !strings.Contains(err.Error(), "device_token=REDACTED") {
		t.Errorf("device token should be redacted, got %v", err)
	}
}

func TestRedactURL(t *testing.T) {
	tests := map[string]string{
		"https://host/api/x?secret=s&title_id=1":   "https://host/api/x?secret=REDACTED&title_id=1",
		"https://host/api/register?device_token=t": "https://host/api/register?device_token=REDACTED",
		"https://host/api/x?title_id=1":            "https://host/api/x?title_id=1",
		"https://host/%zz?secret=s":                "https://host/%zz",
	}
	for in, want := range tests {
		if got := redactURL(in); got != want {
			t.Errorf("redactURL(%q): wanted %q, got %q", in, want, got)
		}
	}
}
