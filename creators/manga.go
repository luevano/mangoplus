package creators

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const MangaListPath = "search/titles"

// MangaService: Provides Manga services provided by the API.
type MangaService service

// TitlesDTO: Container for titles (mangas) lists.
type TitlesDTO struct {
	TitleList  *[]Title    `json:"titleList"`
	Pagination *Pagination `json:"pagination"`
}

// Title: Details of the manga itself.
//
// Avoided some repetitive fields or those that are mostly null always.
type Title struct {
	TitleID                 string `json:"titleId"`
	Title                   string `json:"title"`
	Description             string `json:"description"`
	FirstPublishDate        int64  `json:"firstPublishDate"`
	PublishDate             int64  `json:"publishDate"`
	LatestEpisodeContentsID string `json:"latestEpisodeContentsId"`
	LatestEpisodeTitle      string `json:"latestEpisodeTitle"`
	LatestEpisodeNumbering  int    `json:"latestEpisodeNumbering"`
	Locale                  string `json:"locale"`
	UserID                  int    `json:"userId"`
	HandleName              string `json:"handleName"`
	ProfileImageUrl         string `json:"profileImageUrl"`
	ThumbnailURL            string `json:"thumbnailUrl"`
	ReadCount               int    `json:"readCount"`
	LikeCount               int    `json:"likeCount"`
	FavoriteCount           int    `json:"favoriteCount"`
	CommentCount            int    `json:"commentCount"`
}

// List: Get manga list.
//
// The API itself defaults to the available language (usually english)
// when the requested language doesn't exist.
func (s *MangaService) List(query, language string, page int) (TitlesDTO, error) {
	u, _ := url.Parse(BaseAPI)
	u = u.JoinPath(MangaListPath)
	p := map[string]string{
		"keyword":  query,
		"page":     strconv.Itoa(page),
		"pageSize": "30",
		"sort":     "newly",
		"lang":     language,
	}

	urlR, _ := url.Parse(OriginURL)
	urlR = urlR.JoinPath("keywords")
	paramsR := url.Values{}
	paramsR.Set("q", query)
	urlR.RawQuery = paramsR.Encode()
	h := map[string]string{
		"Referer":          urlR.String(),
		"X-Requested-With": "XMLHttpRequest",
	}

	res, err := s.client.Request(context.Background(), http.MethodGet, *u, p, h, nil)
	if err != nil {
		return TitlesDTO{}, err
	}
	titles := res.TitlesDTO
	if titles == nil || titles.TitleList == nil || len(*titles.TitleList) == 0 {
		return TitlesDTO{}, fmt.Errorf("no mangas found for query %q, language %q, page %d", query, language, page)
	}
	return *titles, nil
}
