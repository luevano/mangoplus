// Package mangoplus provides an API wrapper for MangaPlus API.
package mangoplus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	BaseAPI    = "https://jumpg-api.tokyo-cdn.com/api"
	BaseWebAPI = "https://jumpg-webapi.tokyo-cdn.com/api" // Limited to what's available on the browser
	OriginURL  = "https://mangaplus.shueisha.co.jp"       // Not used
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
	IsFeaturedUpdated *bool              `json:"isFeaturedUpdated"`
	RegisterationData *RegisterationData `json:"registerationData"`
	TitleDetailView   *TitleDetailView   `json:"titleDetailView"`
	MangaViewer       *MangaViewer       `json:"mangaViewer"`
	AllTitlesViewV2   *AllTitlesViewV2   `json:"allTitlesViewV2"`
	Languages         *Languages         `json:"languages"`
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
	client  *http.Client
	secret  *string
	options Options

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
//
// Options must be non-nil. Use DefaultOptions for defaults, all fields must be non-empty.
func NewPlusClient(options Options) *PlusClient {
	options.validate()

	client := http.Client{}
	plus := &PlusClient{client: &client, options: options}
	plus.register()

	plus.common.client = plus
	// Reuse the common client for the other services
	plus.Manga = (*MangaService)(&plus.common)
	plus.Page = (*PageService)(&plus.common)

	return plus
}

// Request: Sends a request to the MangaPlus API and decodes into a PlusResponse.
func (c *PlusClient) Request(
	ctx context.Context,
	method string,
	url url.URL,
	params map[string]string,
	headers map[string]string,
	body io.Reader,
) (PlusResponse, error) {
	p := c.getFinalParams(params)
	url.RawQuery = p.Encode()

	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return PlusResponse{}, nil
	}
	req.Header = c.getFinalHeaders(headers)

	resp, err := c.client.Do(req)
	if err != nil {
		return PlusResponse{}, nil
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		var res PlusResponse
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			return PlusResponse{}, fmt.Errorf("failed to decode MangaPlus response: %s", err.Error())
		}
		if res.Error != nil {
			return PlusResponse{}, fmt.Errorf("MangaPlus API error: %s", res.Error.GetErrors())
		}
		if res.Success == nil {
			return PlusResponse{}, fmt.Errorf("MangaPlus API unexpected error (Error and Success fields empty)")
		}
		return res, nil
	default:
		return PlusResponse{}, fmt.Errorf("MangaPlus API status code %d", resp.StatusCode)
	}
}

// build and get the final params used for the request
func (c *PlusClient) getFinalParams(params map[string]string) url.Values {
	p := url.Values{}
	p.Set("os", "android")
	p.Set("os_ver", c.options.OSVersion)
	p.Set("app_ver", c.options.AppVersion)
	p.Set("format", "json")
	if c.secret != nil {
		p.Set("secret", *c.secret)
	}
	for k, v := range params {
		p.Set(k, v)
	}
	return p
}

func (c *PlusClient) getFinalHeaders(headers map[string]string) http.Header {
	h := http.Header{}
	h.Set("Accept", "*/*") // needed?
	h.Set("User-Agent", c.options.UserAgent)
	for k, v := range headers {
		h.Set(k, v)
	}
	return h
}
