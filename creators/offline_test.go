package creators

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Tests that don't hit the API.

// fakeClient: A client whose requests to the API go to handler instead.
func fakeClient(t *testing.T, handler http.HandlerFunc) *CreatorsClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	target, _ := url.Parse(server.URL)

	o := DefaultOptions()
	o.HTTPClient = &http.Client{Transport: roundTripper(func(r *http.Request) (*http.Response, error) {
		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
		return http.DefaultTransport.RoundTrip(r)
	})}
	c, err := NewCreatorsClient(o)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

type roundTripper func(*http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func answer(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}
}

func TestMangaListEmpty(t *testing.T) {
	for _, body := range []string{
		`{"result":"OK","mpcTitlesDto":{"titleList":[],"pagination":{"page":1,"maxPage":0}}}`,
		`{"result":"OK"}`,
	} {
		c := fakeClient(t, answer(body))
		titles, err := c.Manga.List("nothing", "en", 1)
		if err != nil {
			t.Errorf("%s: wanted no results, got %v", body, err)
		}
		if titles.TitleList != nil && len(*titles.TitleList) != 0 {
			t.Errorf("%s: wanted no titles, got %v", body, *titles.TitleList)
		}
	}
}

func TestMangaList(t *testing.T) {
	c := fakeClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/mpc/search/titles" || r.URL.Query().Get("keyword") != "rauch" {
			t.Errorf("unexpected request %s", r.URL)
		}
		answer(`{"result":"OK","mpcTitlesDto":{"titleList":[{"titleId":"fm1","title":"Rauch"}],"pagination":{"page":1,"maxPage":2}}}`)(w, r)
	})
	titles, err := c.Manga.List("rauch", "en", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(*titles.TitleList) != 1 || (*titles.TitleList)[0].TitleID != "fm1" || !titles.Pagination.HasNextPage() {
		t.Errorf("got %+v", titles)
	}
}

func TestChapterListEmpty(t *testing.T) {
	c := fakeClient(t, answer(`{"result":"OK","mpcEpisodesDto":{"episodeList":[],"pagination":{"page":1,"maxPage":0}}}`))
	episodes, err := c.Chapter.List("fm1", 1)
	if err != nil {
		t.Errorf("wanted no chapters, got %v", err)
	}
	if episodes.EpisodeList != nil && len(*episodes.EpisodeList) != 0 {
		t.Errorf("wanted no chapters, got %v", *episodes.EpisodeList)
	}
}

func TestPageGetOffline(t *testing.T) {
	c := fakeClient(t, answer(`{"result":"OK","pageList":[{"pageNo":1,"publicBgImage":"https://img/1.jpg"}]}`))
	pages, err := c.Page.Get("ep1")
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 1 || pages[0].PublicBGImage != "https://img/1.jpg" {
		t.Errorf("got %+v", pages)
	}
}

func TestNotFound(t *testing.T) {
	c := fakeClient(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	if _, err := c.Page.Get("ep1"); !errors.Is(err, ErrNotFound) {
		t.Errorf("wanted ErrNotFound, got %v", err)
	}
}

func TestResultNotOK(t *testing.T) {
	c := fakeClient(t, answer(`{"result":"NG","messages":["bad request"]}`))
	if _, err := c.Page.Get("ep1"); err == nil || !strings.Contains(err.Error(), "bad request") {
		t.Errorf("wanted the api message, got %v", err)
	}
}

// AWS WAF answers with a JavaScript challenge (202) or a captcha (405) only a browser can solve,
// or blocks the request (403).
func TestWAF(t *testing.T) {
	waf := func(status int, action string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if action != "" {
				w.Header().Set("x-amzn-waf-action", action)
			}
			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			w.WriteHeader(status)
			w.Write([]byte(`<!DOCTYPE html><html><head><script src="https://x.token.awswaf.com/challenge.js"></script></head><body></body></html>`))
		}
	}
	for _, tt := range []struct {
		status int
		action string
	}{
		{http.StatusAccepted, "challenge"},
		{http.StatusMethodNotAllowed, "captcha"},
		{http.StatusAccepted, ""},
		{http.StatusForbidden, ""},
	} {
		c := fakeClient(t, waf(tt.status, tt.action))
		_, err := c.Manga.List("rauch", "en", 1)
		if !errors.Is(err, ErrBlocked) || !strings.Contains(err.Error(), "blocked by AWS WAF") {
			t.Errorf("%d %q: wanted ErrBlocked, got %v", tt.status, tt.action, err)
		}
	}
}
