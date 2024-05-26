package creators

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const PageListPath = "episodes/pageList/%s/"

// PageService: Provides Page services provided by the API.
type PageService service

// Page: Details for each page of a chapter.
//
// Avoided some repetitive fields or those that are mostly null always.
type Page struct {
	PageNo        int    `json:"pageNo"`
	PublicBGImage string `json:"publicBgImage"`
	BGHeight      int    `json:"bgHeight"`
	BGWidth       int    `json:"bgWidth"`
}

// Get: Get list of all chapter pages.
func (s *PageService) Get(id string) ([]Page, error) {
	u, _ := url.Parse(BaseAPI)
	u = u.JoinPath(fmt.Sprintf(PageListPath, id))

	urlR, _ := url.Parse(OriginURL)
	urlR = urlR.JoinPath(fmt.Sprintf("episodes/%s", id))
	h := map[string]string{
		"Referer":          urlR.String(),
		"X-Requested-With": "XMLHttpRequest",
	}

	res, err := s.client.Request(context.Background(), http.MethodGet, *u, nil, h, nil)
	if err != nil {
		return nil, err
	}

	pages := res.PageList
	if pages == nil {
		return nil, fmt.Errorf("no pages found for chapter id %q", id)
	}
	return *pages, nil
}
