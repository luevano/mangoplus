package creators

import (
	"os"
	"testing"
)

// Tests against the live API, only run with MANGOPLUS_LIVE=1 to avoid hitting the API on every run.

var client = NewCreatorsClient(DefaultOptions())

func skipUnlessLive(t *testing.T) {
	if os.Getenv("MANGOPLUS_LIVE") != "1" {
		t.Skip("set MANGOPLUS_LIVE=1 to run tests against the live API")
	}
}

//
// manga.go
//

func TestMangaGet(t *testing.T) {
	skipUnlessLive(t)
	// Rauch
	_, err := client.Manga.List("Rauch", "en", 1)
	if err != nil {
		t.Errorf("Getting manga for query: %s\n", err.Error())
	}
}

//
// chapter.go
//

func TestChapterGet(t *testing.T) {
	skipUnlessLive(t)
	// Rauch chapter list
	_, err := client.Chapter.List("fm2304211111002600024871222", 1)
	if err != nil {
		t.Errorf("Getting chapter list for id: %s\n", err.Error())
	}
}

//
// page.go
//

func TestPageGet(t *testing.T) {
	skipUnlessLive(t)
	// Rauch first chapter pages
	_, err := client.Page.Get("7d2304211111002640024871222")
	if err != nil {
		t.Errorf("Getting chapter pages for id: %s\n", err.Error())
	}
}
