package mangoplus

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	AllMangaPath = "title_list/allV2"
	MangaPath    = "title_detailV3"
)

// MangaService: Provides Manga services provided by the API.
type MangaService service

type TitleDetailView struct {
	Title                    Title              `json:"title"`
	TitleImageUrl            string             `json:"titleImageUrl"`
	Overview                 string             `json:"overview"`
	NextTimeStamp            int                `json:"nextTimeStamp"`
	ViewingPeriodDescription string             `json:"viewingPeriodDescription"`
	ChapterListGroup         []ChapterListGroup `json:"chapterListGroup"`
	IsSimulReleased          bool               `json:"isSimulReleased"`
	Rating                   Rating             `json:"rating"`
	NumberOfViews            int                `json:"numberOfViews"`
	RegionCode               string             `json:"regionCode"`
	TitleLabels              TitleLabels        `json:"titleLabels"`
	Label                    *Label             `json:"label"`
	IsFirstTimeFree          bool               `json:"isFirstTimeFree"`
}

type TitleLabels struct {
	ReleaseSchedule ReleaseSchedule `json:"releaseSchedule"`
	IsSimulpub      bool            `json:"isSimulpub"`
	PlanType        string          `json:"planType"`
}

type Label struct {
	Label LabelCode `json:"label"`
}

type AllTitlesViewV2 struct {
	AllTitlesGroup []AllTitlesGroup `json:"allTitlesGroup"`
}

type AllTitlesGroup struct {
	TheTitle string  `json:"theTitle"`
	Titles   []Title `json:"titles"`
}

type Title struct {
	TitleID           int       `json:"titleId"`
	Name              string    `json:"name"`
	Author            string    `json:"author"`
	PortraitImageURL  string    `json:"portraitImageUrl"`
	Language          *Language `json:"language"`
	ViewCount         *int      `json:"viewCount"`
	TitleUpdateStatus *string   `json:"titleUpdateStatus"`
}

// TODO: make similar helper methods

// GetTitle: Get title of the manga.
// func (m *Manga) GetTitle(langCode string) string {
// 	if title := m.Attributes.Title.GetLocalString(langCode); title != "" {
// 		return title
// 	}
// 	return m.Attributes.AltTitles.GetLocalString(langCode)
// }

// GetDescription: Get description of the manga.
// func (m *Manga) GetDescription(langCode string) string {
// 	return m.Attributes.Description.GetLocalString(langCode)
// }

// Get: Get manga details by ID.
func (s *MangaService) Get(id string) (TitleDetailView, error) {
	u, _ := url.Parse(BaseAPI)
	u = u.JoinPath(MangaPath)

	p := url.Values{}
	p.Set("title_id", id)
	p.Set("format", "json")

	u.RawQuery = p.Encode()

	res, err := s.client.Request(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return TitleDetailView{}, err
	}
	titleDetail := res.Success.TitleDetailView
	if titleDetail == nil {
		return TitleDetailView{}, fmt.Errorf("Error: no details for manga id %s", id)
	}
	return *titleDetail, nil
}

// All: Get list of all manga.
func (s *MangaService) All() ([]AllTitlesGroup, error) {
	u, _ := url.Parse(BaseAPI)
	u = u.JoinPath(AllMangaPath)

	p := url.Values{}
	p.Set("format", "json")

	u.RawQuery = p.Encode()

	res, err := s.client.Request(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	allTitles := res.Success.AllTitlesViewV2
	if allTitles != nil {
		return allTitles.AllTitlesGroup, nil
	}
	return nil, fmt.Errorf("Error: no titles found")
}
