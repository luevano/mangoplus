package mangoplus

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"
)

// Web api responses saved in testdata/web_*.bin (One Piece, title 100020 / chapter 1000486,
// no secret, 2026-10-16), trimmed.

var uuidRegexp = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func newWebTestClient(t *testing.T, options Options) *PlusClient {
	t.Helper()
	files := map[string]string{
		"/" + AllMangaWebPath: "web_allV2.bin",
		"/" + MangaPath:       "web_title_detailV3.bin",
		"/" + PagePathApp:     "web_manga_viewer_v3.bin",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		for _, param := range []string{"secret", "os", "os_ver", "app_ver"} {
			if q.Has(param) {
				t.Errorf("%s: web api request sent %q", r.URL.Path, param)
			}
		}
		if got := r.Header.Get("Session-Token"); !uuidRegexp.MatchString(got) {
			t.Errorf("%s: session token %q isn't a uuid", r.URL.Path, got)
		}
		if got := r.Header.Get("Origin"); got != webOrigin {
			t.Errorf("%s: origin %q", r.URL.Path, got)
		}
		if got := r.Header.Get("User-Agent"); got != defaultWebUserAgent {
			t.Errorf("%s: user agent %q", r.URL.Path, got)
		}
		file, ok := files[r.URL.Path]
		if !ok {
			t.Errorf("unexpected web api request %s", r.URL.Path)
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

	c, err := NewPlusClient(options)
	if err != nil {
		t.Fatal(err)
	}
	if c.baseAPI != BaseWebAPI {
		t.Fatalf("base api %q, want the web api", c.baseAPI)
	}
	if c.options.RequestInterval < time.Second {
		t.Errorf("request interval %s, the web api needs 1s", c.options.RequestInterval)
	}
	c.baseAPI = server.URL
	c.options.RequestInterval = 0
	return c
}

func webOptions() Options {
	o := DefaultOptions()
	o.API = APIWeb
	return o
}

func TestWebAPIChoice(t *testing.T) {
	tests := []struct {
		name    string
		api     API
		secret  string
		free    string
		wantWeb bool
	}{
		{"default", "", "", "", false},
		{"app", APIApp, "", "", false},
		{"web", APIWeb, "", "", true},
		{"web ignores secrets", APIWeb, "paid", "free", true},
		{"auto without secrets", APIAuto, "", "", true},
		{"auto with free secret", APIAuto, "", "free", false},
		{"auto with secret", APIAuto, "paid", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := DefaultOptions()
			o.API, o.Secret, o.FreeSecret = tt.api, tt.secret, tt.free
			c, err := NewPlusClient(o)
			if err != nil {
				t.Fatal(err)
			}
			if c.WebAPI() != tt.wantWeb {
				t.Errorf("WebAPI() = %v, want %v", c.WebAPI(), tt.wantWeb)
			}
			if tt.wantWeb && c.secret != "" {
				t.Errorf("web api client kept a secret")
			}
		})
	}

	o := DefaultOptions()
	o.API = "browser"
	if _, err := NewPlusClient(o); err == nil {
		t.Error("unknown API should be an error")
	}
}

func TestWebMangaAll(t *testing.T) {
	c := newWebTestClient(t, webOptions())
	titles, err := c.Manga.All(Language_ENGLISH)
	if err != nil {
		t.Fatal(err)
	}
	if len(titles) != 6 {
		t.Errorf("got %d english titles, want 6", len(titles))
	}
	for _, title := range titles {
		if title.Language != Language_ENGLISH {
			t.Errorf("%q is in %s", title.Name, title.Language)
		}
	}

	spanish, err := c.Manga.All(Language_SPANISH)
	if err != nil {
		t.Fatal(err)
	}
	if len(spanish) != 1 || spanish[0].TitleId != 200108 {
		t.Errorf("spanish titles %v, want only 200108", spanish)
	}
}

func TestWebMangaGet(t *testing.T) {
	c := newWebTestClient(t, webOptions())
	detail, err := c.Manga.Get("100020")
	if err != nil {
		t.Fatal(err)
	}
	if detail.Title.Name != "One Piece" {
		t.Errorf("title %q", detail.Title.Name)
	}
	// the mid lists (#004, #1192) are hidden on the website
	want := []int32{1000486, 1000487, 1000488, 1030380, 1030381, 1030615}
	if len(detail.ChapterListV2) != len(want) {
		t.Fatalf("got %d chapters, want %d", len(detail.ChapterListV2), len(want))
	}
	for i, chapter := range detail.ChapterListV2 {
		if chapter.ChapterId != want[i] {
			t.Errorf("chapter %d is %d, want %d", i, chapter.ChapterId, want[i])
		}
	}
	if detail.TitleLabels.GetReleaseSchedule() == ReleaseSchedule_DISABLED {
		t.Error("no release schedule")
	}
}

func TestWebPageGet(t *testing.T) {
	c := newWebTestClient(t, webOptions())
	pages, err := c.Page.Get("1000486", true, ImageQualitySuperHigh)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 3 {
		t.Fatalf("got %d pages, want 3", len(pages))
	}
	if pages[0].EncryptionKey == "" {
		t.Error("page without encryption key")
	}
	headers := c.ImageHeaders(pages[0].ImageUrl)
	if got := headers["Plus-Vw-Token"]; got != "bcde7e4ac84894354b4234bfcaaddc0b" {
		t.Errorf("Plus-Vw-Token %q", got)
	}
	if got := headers["User-Agent"]; got != defaultWebUserAgent {
		t.Errorf("image user agent %q", got)
	}
}

func TestWebUserUnsupported(t *testing.T) {
	c := newWebTestClient(t, webOptions())
	if _, err := c.User.Subscription(); !errors.Is(err, ErrWebAPIUnsupported) {
		t.Errorf("subscription: %v, want ErrWebAPIUnsupported", err)
	}
	if _, err := c.User.Profile(); !errors.Is(err, ErrWebAPIUnsupported) {
		t.Errorf("profile: %v, want ErrWebAPIUnsupported", err)
	}
}
