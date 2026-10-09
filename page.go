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

	var pages []*MangaPage
	for _, page := range mangaViewer.Pages {
		if page.MangaPage != nil {
			pages = append(pages, page.MangaPage)
		}
	}
	return pages, nil
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
