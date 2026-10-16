package mangoplus

import (
	"context"
	"fmt"
	"net/http"
)

const (
	PagePath    = "manga_viewer"
	PagePathApp = "manga_viewer_v3"
)

type PageService service

// Get: Get list of all chapter pages, using the endpoint set by Options.Viewer.
func (s *PageService) Get(id string, splitImages bool, imageQuality ImageQuality) ([]*MangaPage, error) {
	return s.GetContext(context.Background(), id, splitImages, imageQuality)
}

// GetContext: Same as Get, with a context.
func (s *PageService) GetContext(ctx context.Context, id string, splitImages bool, imageQuality ImageQuality) ([]*MangaPage, error) {
	p := map[string]string{
		"chapter_id":  id,
		"split":       yesNo(splitImages),
		"img_quality": string(imageQuality),
	}

	if s.client.web {
		return s.getWeb(ctx, id, p)
	}

	subscription := s.client.options.Secret != ""
	path := PagePath
	switch s.client.options.Viewer {
	case ViewerLegacy:
		if subscription {
			p["subscription_reading"] = "yes"
		}
	case ViewerApp:
		path = PagePathApp
		p["viewer_mode"] = "horizontal"
		p["ticket_reading"] = "no"
		// The app uses the free read on free tier, the subscription otherwise
		p["free_reading"] = yesNo(!subscription)
		p["subscription_reading"] = yesNo(subscription)
	}

	res, err := s.client.Request(ctx, http.MethodGet, path, p)
	if err != nil {
		return nil, fmt.Errorf("get chapter pages %q: %w", id, err)
	}
	mangaViewer := res.GetMangaViewer()
	if mangaViewer == nil {
		return nil, fmt.Errorf("get chapter pages %q: no manga viewer in response", id)
	}
	return mangaPages(mangaViewer), nil
}

// getWeb: The web api's viewer, its token goes with the images (ImageHeaders).
func (s *PageService) getWeb(ctx context.Context, id string, p map[string]string) ([]*MangaPage, error) {
	p["clang"] = webContentLanguage
	res, err := s.client.Request(ctx, http.MethodGet, PagePathApp, p)
	if err != nil {
		return nil, fmt.Errorf("get chapter pages %q: %w", id, err)
	}
	mangaViewer := res.GetMangaViewer()
	if mangaViewer == nil {
		return nil, fmt.Errorf("get chapter pages %q: no manga viewer in response", id)
	}
	s.client.viewTokenMu.Lock()
	s.client.viewToken = mangaViewer.ViewToken
	s.client.viewTokenMu.Unlock()
	return mangaPages(mangaViewer), nil
}

func mangaPages(mangaViewer *MangaViewer) []*MangaPage {
	var pages []*MangaPage
	for _, page := range mangaViewer.Pages {
		if page.MangaPage != nil {
			pages = append(pages, page.MangaPage)
		}
	}
	return pages
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
