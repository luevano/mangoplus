package creators

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// There are more fields in the responses, but I found most of them
// being null on all my testings, could update if needed.
// Some of the fields included are not required but they're included anyways.

const (
	BaseAPI   = "https://medibang.com/api/mpc"
	OriginURL = "https://medibang.com/mpc"
)

// CreatorsResponse: Generic MangaPlusCreators API response type, most responses have this structure.
type CreatorsResponse struct {
	Result      string       `json:"result"`
	Messages    *[]string    `json:"messages"`
	TitlesDTO   *TitlesDTO   `json:"mpcTitlesDto"`
	EpisodesDTO *EpisodesDTO `json:"mpcEpisodesDto"`
	LayoutType  *string      `json:"layoutType"`
	PageList    *[]Page      `json:"pageList"`
}

func (r CreatorsResponse) GetMessages() string {
	if r.Messages == nil {
		return "no messages"
	}
	return strings.Join(*r.Messages, ", ")
}

// Pagination: Information on available search pages.
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Count    int `json:"count"`
	MaxPage  int `json:"maxPage"`
	Start    int `json:"start"`
}

func (p Pagination) HasNextPage() bool {
	return p.Page < p.MaxPage
}

// CreatorsClient: The MangaPlusCreators client.
type CreatorsClient struct {
	client  *http.Client
	options Options

	common service

	// Services for MangaPlusCreators API.
	Manga   *MangaService
	Chapter *ChapterService
	Page    *PageService
}

// service: Wrapper for CreatorsClient.
type service struct {
	client *CreatorsClient
}

// NewCreatorsClient: New MangaPlusCreators client.
//
// Options must be non-nil. Use DefaultOptions for defaults, all fields must be non-empty.
func NewCreatorsClient(options Options) (*CreatorsClient, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("mangapluscreators: invalid options: %w", err)
	}

	client := http.Client{}
	creators := &CreatorsClient{client: &client, options: options}

	creators.common.client = creators
	// Reuse the common client for the other services
	creators.Manga = (*MangaService)(&creators.common)
	creators.Chapter = (*ChapterService)(&creators.common)
	creators.Page = (*PageService)(&creators.common)

	return creators, nil
}

// Request: Sends a request to the MangaPlus API and decodes into a CreatorsResponse.
func (c *CreatorsClient) Request(
	ctx context.Context,
	method string,
	url url.URL,
	params map[string]string,
	headers map[string]string,
	body io.Reader,
) (CreatorsResponse, error) {
	p := c.getFinalParams(params)
	url.RawQuery = p.Encode()
	// prefix for errors that aren't from the API itself
	op := fmt.Sprintf("mangapluscreators: %s %s", method, url.Path)

	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return CreatorsResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	req.Header = c.getFinalHeaders(headers)

	resp, err := c.client.Do(req)
	if err != nil {
		return CreatorsResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	// AWS WAF answers 202 with a challenge only a browser can solve
	if resp.Header.Get("x-amzn-waf-action") != "" {
		return CreatorsResponse{}, fmt.Errorf("%s: blocked by an AWS WAF bot challenge", op)
	}
	if resp.StatusCode != http.StatusOK {
		return CreatorsResponse{}, &HTTPError{Method: method, Path: url.Path, StatusCode: resp.StatusCode}
	}

	var res CreatorsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return CreatorsResponse{}, fmt.Errorf("%s: decode response: %w", op, err)
	}
	if res.Result != "OK" {
		return CreatorsResponse{}, fmt.Errorf("%s: result %q: %s", op, res.Result, res.GetMessages())
	}
	return res, nil
}

// build and get the final params used for the request
func (c *CreatorsClient) getFinalParams(params map[string]string) url.Values {
	p := url.Values{}
	p.Set("_", strconv.FormatInt(time.Now().UnixMilli(), 10))
	for k, v := range params {
		p.Set(k, v)
	}
	return p
}

func (c *CreatorsClient) getFinalHeaders(headers map[string]string) http.Header {
	h := http.Header{}
	h.Set("User-Agent", c.options.UserAgent)
	for k, v := range headers {
		h.Set(k, v)
	}
	return h
}
