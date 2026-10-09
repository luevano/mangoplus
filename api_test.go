package mangoplus

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
)

// Offline tests, the server answers with real responses saved in testdata/
// (Ghost Fixers, title 100310 / chapter 1020497, from a free tier device).

const testSecret = "0123456789abcdef0123456789abcdef"

// newTestClient: Client pointing to a local server, which records the query of the last request.
func newTestClient(t *testing.T, options Options) (*PlusClient, *url.Values) {
	t.Helper()
	var lastQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastQuery = r.URL.Query()

		// built here, so no real device secret is stored in testdata
		if r.URL.Path == "/register" {
			res := &Response{Success: &SuccessResult{RegistrationData: &RegistrationData{DeviceSecret: testSecret}}}
			body, _ := proto.Marshal(res)
			w.Write(body)
			return
		}

		files := map[string]string{
			"/title_detailV3":     "title_detailV3.bin",
			"/title_list/all_v3":  "all_v3.bin",
			"/manga_viewer":       "manga_viewer_v3.bin",
			"/manga_viewer_v3":    "manga_viewer_v3.bin",
			"/profile":            "profile.bin",
			"/subscription":       "subscription.bin",
			"/manga_viewer_error": "error.bin",
		}
		file, ok := files[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		body, err := os.ReadFile("testdata/" + file)
		if err != nil {
			t.Errorf("reading testdata: %s", err)
		}
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	c := NewPlusClient(options)
	c.baseAPI = server.URL
	return c, &lastQuery
}

func TestRegister(t *testing.T) {
	c, lastQuery := newTestClient(t, DefaultOptions())

	if _, err := c.User.Profile(); err != nil {
		t.Fatal(err)
	}
	if got := lastQuery.Get("secret"); got != testSecret {
		t.Errorf("request after registering should use the new secret, got %q", got)
	}
}

func TestSecretSkipsRegister(t *testing.T) {
	o := DefaultOptions()
	o.Secret = "my-secret"
	c, lastQuery := newTestClient(t, o)

	if _, err := c.User.Profile(); err != nil {
		t.Fatal(err)
	}
	if got := lastQuery.Get("secret"); got != "my-secret" {
		t.Errorf("wanted the given secret, got %q", got)
	}
}

func TestOnRegister(t *testing.T) {
	var registered string
	o := DefaultOptions()
	o.OnRegister = func(secret string) { registered = secret }
	c, _ := newTestClient(t, o)

	if _, err := c.User.Profile(); err != nil {
		t.Fatal(err)
	}
	if registered != testSecret {
		t.Errorf("OnRegister should get the new secret, got %q", registered)
	}
}

func TestFreeSecret(t *testing.T) {
	o := DefaultOptions()
	o.FreeSecret = "free-secret"
	o.OnRegister = func(string) { t.Error("shouldn't register a new device") }
	c, lastQuery := newTestClient(t, o)

	if _, err := c.Page.Get("1020497", false, ImageQualityHigh); err != nil {
		t.Fatal(err)
	}
	if got := lastQuery.Get("secret"); got != "free-secret" {
		t.Errorf("wanted the free secret, got %q", got)
	}
	if got := lastQuery.Get("subscription_reading"); got != "no" {
		t.Errorf("free secret shouldn't use the subscription, got subscription_reading %q", got)
	}

	// the paid one takes precedence
	o.Secret = "paid-secret"
	c, lastQuery = newTestClient(t, o)
	if _, err := c.User.Profile(); err != nil {
		t.Fatal(err)
	}
	if got := lastQuery.Get("secret"); got != "paid-secret" {
		t.Errorf("wanted the paid secret, got %q", got)
	}
}

func TestMangaGet(t *testing.T) {
	c, lastQuery := newTestClient(t, DefaultOptions())

	titleDetail, err := c.Manga.Get("100310")
	if err != nil {
		t.Fatal(err)
	}
	if got := lastQuery.Get("title_id"); got != "100310" {
		t.Errorf("wrong title_id param: %q", got)
	}
	if name := titleDetail.GetTitle().GetName(); name != "GHOST FIXERS" {
		t.Errorf("wanted %q, got %q", "GHOST FIXERS", name)
	}
	if len(titleDetail.ChapterListV2) == 0 {
		t.Fatal("manga has no chapters")
	}
	if name := titleDetail.ChapterListV2[0].Name; name != "#001" {
		t.Errorf("wanted first chapter %q, got %q", "#001", name)
	}
	if label := titleDetail.GetLabel().GetLabel(); label.Long() == "" {
		t.Errorf("label %s has no long name", label)
	}
}

func TestMangaAll(t *testing.T) {
	c, lastQuery := newTestClient(t, DefaultOptions())

	titles, err := c.Manga.All(Language_SPANISH)
	if err != nil {
		t.Fatal(err)
	}
	if got := lastQuery.Get("clang"); got != "esp" {
		t.Errorf("wrong clang param: %q", got)
	}
	for _, title := range titles {
		if title.Language != Language_SPANISH {
			t.Fatalf("manga %q is not in spanish: %s", title.Name, title.Language)
		}
	}
}

func TestPageGet(t *testing.T) {
	tests := []struct {
		viewer       Viewer
		secret       string
		subscription string
	}{
		{ViewerApp, "", "no"},
		{ViewerApp, "my-secret", "yes"},
		{ViewerLegacy, "", ""},
		{ViewerLegacy, "my-secret", "yes"},
	}
	for _, tt := range tests {
		o := DefaultOptions()
		o.Viewer = tt.viewer
		o.Secret = tt.secret
		c, lastQuery := newTestClient(t, o)

		pages, err := c.Page.Get("1020497", false, ImageQualityHigh)
		if err != nil {
			t.Fatal(err)
		}
		if len(pages) == 0 || pages[0].ImageUrl == "" {
			t.Errorf("%s viewer: no pages", tt.viewer)
		}
		if got := lastQuery.Get("subscription_reading"); got != tt.subscription {
			t.Errorf("%s viewer, secret %q: wanted subscription_reading %q, got %q", tt.viewer, tt.secret, tt.subscription, got)
		}
	}
}

func TestUser(t *testing.T) {
	c, _ := newTestClient(t, DefaultOptions())

	subscription, err := c.User.Subscription()
	if err != nil {
		t.Fatal(err)
	}
	if subscription.PlanType != "basic" {
		t.Errorf("wanted the basic plan, got %q", subscription.PlanType)
	}
	if _, err := c.User.Profile(); err != nil {
		t.Fatal(err)
	}
}

func TestErrorResponse(t *testing.T) {
	c, _ := newTestClient(t, DefaultOptions())

	_, err := c.Request(context.Background(), http.MethodGet, "manga_viewer_error", nil)
	if err == nil || !strings.Contains(err.Error(), "(11103)") {
		t.Errorf("wanted the server error message, got %v", err)
	}
}

func TestContextCanceled(t *testing.T) {
	c, _ := newTestClient(t, DefaultOptions())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.Manga.GetContext(ctx, "100310"); !errors.Is(err, context.Canceled) {
		t.Errorf("wanted context.Canceled, got %v", err)
	}
}
