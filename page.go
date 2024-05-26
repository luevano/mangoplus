package mangoplus

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	PagePath = "manga_viewer"
)

// PageService: Provides Page services provided by the API.
type PageService service

type MangaViewer struct {
	Pages            []Page    `json:"pages"`
	ChapterID        int       `json:"chapterId"`
	Chapters         []Chapter `json:"chapters"` // Probably not really needed
	TitleName        string    `json:"titleName"`
	ChapterName      string    `json:"chapterName"`
	NumberOfComments int       `json:"numberOfComments"`
	TitleID          int       `json:"titleId"`
	RegionCode       string    `json:"regionCode"`
	TitleLanguage    string    `json:"titleLanguage"` // This uses a completely different format than that of other Language fields
}

type Page struct {
	MangaPage *MangaPage `json:"mangaPage"`
	// There are other pages at the end that are served as ads or other info
}

type MangaPage struct {
	ImageURL      string  `json:"imageUrl"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	Type          *string `json:"type"` // Could be "DOUBLE" (double paged)
	EncryptionKey *string `json:"encryptionKey"`
}

// Get: Get list of all chapter pages.
func (s *PageService) Get(id string, splitImages bool, imageQuality ImageQuality) ([]MangaPage, error) {
	u, _ := url.Parse(BaseAPI)
	u = u.JoinPath(PagePath)

	split := "no"
	if splitImages {
		split = "yes"
	}
	p := map[string]string{
		"chapter_id":  id,
		"split":       split,
		"img_quality": string(imageQuality),
	}

	res, err := s.client.Request(context.Background(), http.MethodGet, *u, p, nil, nil)
	if err != nil {
		return nil, err
	}
	mangaViewer := res.Success.MangaViewer
	if mangaViewer == nil {
		return nil, fmt.Errorf("unexpected error while getting pages for chapter id %q", id)
	}

	var pages []MangaPage
	for _, page := range mangaViewer.Pages {
		if page.MangaPage != nil {
			pages = append(pages, *page.MangaPage)
		}
	}
	return pages, nil
}
