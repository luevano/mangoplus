package creators

import "testing"

var client = NewCreatorsClient(DefaultOptions())

//
// manga.go
//

func TestMangaGet(t *testing.T) {
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
	// Rauch first chapter pages
	_, err := client.Page.Get("7d2304211111002640024871222")
	if err != nil {
		t.Errorf("Getting chapter pages for id: %s\n", err.Error())
	}
}
