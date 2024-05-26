package creators

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const ChapterListPath = "titles/%s/episodes/"

// ChapterService: Provides Manga services provided by the API.
type ChapterService service

// EpisodesDTO: Container for episodes (chapters) lists.
type EpisodesDTO struct {
	EpisodeList *[]Episode  `json:"episodeList"`
	Pagination  *Pagination `json:"pagination"`
}

// Episode: Details for a chapter of the manga.
//
// Avoided some repetitive fields or those that are mostly null always.
type Episode struct {
	EpisodeID         string `json:"episodeId"`
	EpisodeTitle      string `json:"episodeTitle"`
	Afterword         string `json:"afterword"`
	Numbering         int    `json:"numbering"`
	FirstPageIsSpread bool   `json:"firstPageIsSpread"`
	FirstPublishDate  int64  `json:"firstPublishDate"`
	PublishDate       int64  `json:"publishDate"`
	TitleId           string `json:"titleId"`
	ThumbnailUrl      string `json:"thumbnailUrl"`
	ReadCount         int    `json:"readCount"`
	LikeCount         int    `json:"likeCount"`
	FavoriteCount     int    `json:"favoriteCount"`
	CommentCount      int    `json:"commentCount"`
	Oneshot           bool   `json:"oneshot"`
}

// List: Get chapter list.
func (s *ChapterService) List(id string, page int) (EpisodesDTO, error) {
	u, _ := url.Parse(BaseAPI)
	u = u.JoinPath(fmt.Sprintf(ChapterListPath, id))
	p := map[string]string{
		"page":     strconv.Itoa(page),
		"pageSize": "200",
	}

	urlR, _ := url.Parse(OriginURL)
	urlR = urlR.JoinPath(fmt.Sprintf("titles/%s", id))
	h := map[string]string{
		"Referer":          urlR.String(),
		"X-Requested-With": "XMLHttpRequest",
	}

	res, err := s.client.Request(context.Background(), http.MethodGet, *u, p, h, nil)
	if err != nil {
		return EpisodesDTO{}, err
	}

	episodes := res.EpisodesDTO
	if episodes == nil  || episodes.EpisodeList == nil || len(*episodes.EpisodeList) == 0 {
		return EpisodesDTO{}, fmt.Errorf("no chapters found for manga id %q, page %d", id, page)
	}
	return *episodes, nil
}
