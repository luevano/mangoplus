package mangoplus

import (
	"context"
	"fmt"
	"net/http"
)

const (
	AllMangaPath = "title_list/all_v3"
	MangaPath    = "title_detailV3"
)

// The catalog is split in these groups, same as the tabs in the app.
var allMangaTypes = []string{"serializing", "completed", "one-shot"}

type MangaService service

// Get: Get manga details by ID.
func (s *MangaService) Get(id string) (*TitleDetailView, error) {
	return s.GetContext(context.Background(), id)
}

// GetContext: Same as Get, with a context.
func (s *MangaService) GetContext(ctx context.Context, id string) (*TitleDetailView, error) {
	p := map[string]string{"title_id": id}

	res, err := s.client.Request(ctx, http.MethodGet, MangaPath, p)
	if err != nil {
		return nil, err
	}
	titleDetail := res.GetTitleDetailView()
	if titleDetail == nil {
		return nil, fmt.Errorf("no details for manga id %q", id)
	}
	return titleDetail, nil
}

// All: Get list of all manga in the given language.
func (s *MangaService) All(language Language) ([]*Title, error) {
	return s.AllContext(context.Background(), language)
}

// AllContext: Same as All, with a context.
func (s *MangaService) AllContext(ctx context.Context, language Language) ([]*Title, error) {
	var titles []*Title
	for _, mangaType := range allMangaTypes {
		p := map[string]string{
			"type":  mangaType,
			"lang":  language.APICode(),
			"clang": language.APICode(),
		}

		res, err := s.client.Request(ctx, http.MethodGet, AllMangaPath, p)
		if err != nil {
			return nil, err
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
		return nil, fmt.Errorf("no mangas found")
	}
	return titles, nil
}
