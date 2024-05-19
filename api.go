// Package mangoplus provides an API wrapper for MangaPlus API.
package mangoplus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

const (
	OriginURL   = "https://mangaplus.shueisha.co.jp"
	BaseAPI     = "https://jumpg-webapi.tokyo-cdn.com/api"
	DefaultUUID = "710aedb9-8614-4fe2-84d5-90624d5af04d"
)

// PlusResponse: Generic MangaPlus API response type, most responses have this structure.
type PlusResponse struct {
	Success *SuccessResponse `json:"success"`
	Error   *ErrorResponse   `json:"error"`
}

// ErrorResponse: Generic error response.
type ErrorResponse struct {
	// Not sure if English/Spanish are always the ones that
	// appear on top, could be specific to me.
	EnglishPopup *Popup   `json:"englishPopup"`
	SpanishPopup *Popup   `json:"spanishPopup"`
	Popups       *[]Popup `json:"popups"`
}

// SuccessResponse: Generic success response.
type SuccessResponse struct {
	IsFeaturedUpdated *bool            `json:"isFeaturedUpdated"`
	TitleDetailView   *TitleDetailView `json:"titleDetailView"`
	MangaViewer       *MangaViewer     `json:"mangaViewer"`
	AllTitlesViewV2   *AllTitlesViewV2 `json:"allTitlesViewV2"`
	Languages         *Languages       `json:"languages"`
}

// Languages: Part of the response when requesting all of the manga.
//
// Not really used.
type Languages struct {
	DefaultUILanguage         Language `json:"defaultUiLanguage"`
	DefaultContentLanguageOne Language `json:"defaultContentLanguageOne"`
	AvailableLanguages        []struct {
		Language    *Language `json:"language"`
		TitlesCount int       `json:"titlesCount"`
	} `json:"availableLanguages"`
}

// PlusClient: The MangaPlus client.
type PlusClient struct {
	client *http.Client
	header http.Header

	common service

	// Services for MangaPlus API.
	Manga *MangaService
	Page  *PageService
}

// service: Wrapper for PlusClient.
type service struct {
	client *PlusClient
}

// NewPlusClient: New MangaPlus client.
func NewPlusClient() *PlusClient {
	client := http.Client{}
	header := http.Header{}

	// Not sure if these headers are needed.
	// The page downloader uses different client/headers.
	header.Set("Origin", OriginURL)
	header.Set("Referer", fmt.Sprintf("%s/", OriginURL))
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	randUUID, err := uuid.NewRandom()
	if err != nil {
		header.Set("SESSION-TOKEN", DefaultUUID)
	} else {
		header.Set("SESSION-TOKEN", randUUID.String())
	}

	plus := &PlusClient{
		client: &client,
		header: header,
	}
	plus.common.client = plus

	// Reuse the common client for the other services
	plus.Manga = (*MangaService)(&plus.common)
	plus.Page = (*PageService)(&plus.common)

	return plus
}

// Request: Sends a request to the MangaPlus API and decodes into a PlusResponse.
func (c *PlusClient) Request(ctx context.Context, method, url string, body io.Reader) (PlusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return PlusResponse{}, nil
	}
	req.Header = c.header

	resp, err := c.client.Do(req)
	if err != nil {
		return PlusResponse{}, nil
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		var res PlusResponse
		err = json.NewDecoder(resp.Body).Decode(&res)
		if res.Error != nil {
			return PlusResponse{}, fmt.Errorf("Error: %s", res.Error.GetErrors())
		}
		if res.Success == nil {
			return PlusResponse{}, fmt.Errorf("Error: didn't receive neither an error or success response")
		}
		return res, nil
	default:
		return PlusResponse{}, fmt.Errorf("Error: status code %d", resp.StatusCode)
	}
}
