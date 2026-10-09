package mangoplus

import (
	"os"
	"testing"
)

// Tests against the live API, only run with MANGOPLUS_LIVE=1 to avoid hitting the API on every run.

var client = NewPlusClient(DefaultOptions())

func skipUnlessLive(t *testing.T) {
	if os.Getenv("MANGOPLUS_LIVE") != "1" {
		t.Skip("set MANGOPLUS_LIVE=1 to run tests against the live API")
	}
}

//
// manga.go
//

func TestLiveMangaGet(t *testing.T) {
	skipUnlessLive(t)
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

func TestLiveMangaAll(t *testing.T) {
	skipUnlessLive(t)
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

func TestLivePageGet(t *testing.T) {
	skipUnlessLive(t)
	// Ghost fixers' first chapter
	pages, err := client.Page.Get("1020497", false, ImageQualityHigh)
	if err != nil {
		t.Fatalf("Getting chapter pages by id: %s\n", err.Error())
	}
	if len(pages) == 0 {
		t.Errorf("Chapter has no pages")
	}
}

func TestLivePageGetLegacyViewer(t *testing.T) {
	skipUnlessLive(t)
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

func TestLiveUserSubscription(t *testing.T) {
	skipUnlessLive(t)
	subscription, err := client.User.Subscription()
	if err != nil {
		t.Fatalf("Getting subscription: %s\n", err.Error())
	}
	if subscription.PlanType != "basic" {
		t.Errorf("Newly registered device should have the basic plan, got %q", subscription.PlanType)
	}
}

func TestLiveUserProfile(t *testing.T) {
	skipUnlessLive(t)
	if _, err := client.User.Profile(); err != nil {
		t.Fatalf("Getting profile: %s\n", err.Error())
	}
}
