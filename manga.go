package mangoplus

import (
	"context"
	"fmt"
	"net/http"
)

const (
	AllMangaPath = "title_list/all_v3"
	MangaPath    = "title_detailV3"
	// AllMangaWebPath: The web api's list of every title, in one request.
	AllMangaWebPath = "title_list/allV2"
)

// The catalog is split in these groups, same as the tabs in the app.
var allMangaTypes = []string{"serializing", "completed", "one-shot"}

type MangaService service

// Get: Get manga details by ID.
func (s *MangaService) Get(id string) (*TitleDetailView, error) {
	return s.GetContext(context.Background(), id)
}

// GetContext: Same as Get, with a context.
//
// With the web api, ChapterListV2 is filled with the chapters the website shows (the first and
// latest ones of ChapterListGroup).
func (s *MangaService) GetContext(ctx context.Context, id string) (*TitleDetailView, error) {
	p := map[string]string{"title_id": id}
	if s.client.web {
		p["clang"] = webContentLanguage
	}

	res, err := s.client.Request(ctx, http.MethodGet, MangaPath, p)
	if err != nil {
		return nil, fmt.Errorf("get manga %q: %w", id, err)
	}
	titleDetail := res.GetTitleDetailView()
	if titleDetail == nil {
		return nil, fmt.Errorf("get manga %q: no title details in response", id)
	}
	if s.client.web && len(titleDetail.ChapterListV2) == 0 {
		for _, group := range titleDetail.ChapterListGroup {
			titleDetail.ChapterListV2 = append(titleDetail.ChapterListV2, group.FirstChapterList...)
			titleDetail.ChapterListV2 = append(titleDetail.ChapterListV2, group.LastChapterList...)
		}
	}
	return titleDetail, nil
}

// All: Get list of all manga in the given language.
func (s *MangaService) All(language Language) ([]*Title, error) {
	return s.AllContext(context.Background(), language)
}

// AllContext: Same as All, with a context.
func (s *MangaService) AllContext(ctx context.Context, language Language) ([]*Title, error) {
	if s.client.web {
		return s.allWeb(ctx, language)
	}
	var titles []*Title
	for _, mangaType := range allMangaTypes {
		p := map[string]string{
			"type":  mangaType,
			"lang":  language.APICode(),
			"clang": language.APICode(),
		}

		res, err := s.client.Request(ctx, http.MethodGet, AllMangaPath, p)
		if err != nil {
			return nil, fmt.Errorf("list all %s manga: %w", mangaType, err)
		}
		// Each group has the title in all of its languages, regardless of "clang"
		for _, group := range res.GetSearchView().GetAllTitlesGroup() {
			for _, title := range group.Titles {
				if title.Language == language {
					titles = append(titles, title)
				}
			}
		}
	}
	if len(titles) == 0 {
		return nil, fmt.Errorf("list all manga in %s: %w", language, ErrNotFound)
	}
	return titles, nil
}

func (s *MangaService) allWeb(ctx context.Context, language Language) ([]*Title, error) {
	res, err := s.client.Request(ctx, http.MethodGet, AllMangaWebPath, nil)
	if err != nil {
		return nil, fmt.Errorf("list all manga: %w", err)
	}
	var titles []*Title
	for _, group := range res.GetAllTitlesView().GetAllTitlesGroup() {
		for _, title := range group.Titles {
			if title.Language == language && title.TitleId != 0 {
				titles = append(titles, title)
			}
		}
	}
	if len(titles) == 0 {
		return nil, fmt.Errorf("list all manga in %s: %w", language, ErrNotFound)
	}
	return titles, nil
}
