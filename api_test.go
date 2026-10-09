package mangoplus

import "testing"

var client = NewPlusClient(DefaultOptions())

//
// manga.go
//

func TestMangaGet(t *testing.T) {
	// Ghost fixers
	titleDetail, err := client.Manga.Get("100310")
	if err != nil {
		t.Fatalf("Getting manga by id: %s\n", err.Error())
	}
	if name := titleDetail.GetTitle().GetName(); name != "GHOST FIXERS" {
		t.Errorf("Manga name doesn't match expected; wanted %q, got %q", "GHOST FIXERS", name)
	}
	if len(titleDetail.ChapterListV2) == 0 {
		t.Errorf("Manga has no chapters")
	}
}

func TestMangaAll(t *testing.T) {
	titles, err := client.Manga.All(Language_ENGLISH)
	if err != nil {
		t.Fatalf("Getting all manga: %s\n", err.Error())
	}
	for _, title := range titles {
		if title.Language != Language_ENGLISH {
			t.Fatalf("Manga %q is not in english: %s", title.Name, title.Language)
		}
	}
}

//
// page.go
//

func TestPageGet(t *testing.T) {
	// Ghost fixers' first chapter
	pages, err := client.Page.Get("1020497", false, ImageQualityHigh)
	if err != nil {
		t.Fatalf("Getting chapter pages by id: %s\n", err.Error())
	}
	if len(pages) == 0 {
		t.Errorf("Chapter has no pages")
	}
}

func TestPageGetLegacyViewer(t *testing.T) {
	o := DefaultOptions()
	o.Viewer = ViewerLegacy
	appClient := NewPlusClient(o)

	// Ghost fixers' first chapter
	pages, err := appClient.Page.Get("1020497", false, ImageQualityHigh)
	if err != nil {
		t.Fatalf("Getting chapter pages by id: %s\n", err.Error())
	}
	if len(pages) == 0 {
		t.Errorf("Chapter has no pages")
	}
}

//
// user.go
//

func TestUserSubscription(t *testing.T) {
	subscription, err := client.User.Subscription()
	if err != nil {
		t.Fatalf("Getting subscription: %s\n", err.Error())
	}
	if subscription.PlanType != "basic" {
		t.Errorf("Newly registered device should have the basic plan, got %q", subscription.PlanType)
	}
}

func TestUserProfile(t *testing.T) {
	if _, err := client.User.Profile(); err != nil {
		t.Fatalf("Getting profile: %s\n", err.Error())
	}
}
