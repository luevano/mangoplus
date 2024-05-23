package mangoplus

import "testing"

var client = NewPlusClient(DefaultOptions())

//
// manga.go
//

func TestMangaGet(t *testing.T) {
	// Ghost fixers
	_, err := client.Manga.Get("100310")
	if err != nil {
		t.Errorf("Getting manga by id: %s\n", err.Error())
	}
}

func TestMangaAll(t *testing.T) {
	_, err := client.Manga.All()
	if err != nil {
		t.Errorf("Getting all manga: %s\n", err.Error())
	}
}

//
// page.go
//

func TestPageGet(t *testing.T) {
	// Ghost fixers' first chapter
	_, err := client.Page.Get("1020497", false, ImageQualityHigh)
	if err != nil {
		t.Errorf("Getting chapter pages by id: %s\n", err.Error())
	}
}
