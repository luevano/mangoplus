// Package mangoplus provides an API wrapper for MangaPlus API.
package mangoplus

//go:generate go run github.com/bufbuild/buf/cmd/buf@v1.73.0 generate --template buf.gen.yaml

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

const BaseAPI = "https://jumpg-api.tokyo-cdn.com/api"

// PlusClient: The MangaPlus client.
type PlusClient struct {
	client  *http.Client
	options Options
	// BaseAPI, only changed in tests
	baseAPI string

	// Device secret sent on every request, from Options.Secret or registering a new device.
	secret     string
	registerMu sync.Mutex

	lastRequest   time.Time
	lastRequestMu sync.Mutex

	common service

	// Services for MangaPlus API.
	Manga *MangaService
	Page  *PageService
	User  *UserService
}

// service: Wrapper for PlusClient.
type service struct {
	client *PlusClient
}

// NewPlusClient: New MangaPlus client.
//
// Use DefaultOptions for defaults, all fields except Secret must be non-empty.
//
// If no Secret is provided, a new device is registered on the first request.
func NewPlusClient(options Options) *PlusClient {
	if err := options.validate(); err != nil {
		panic(fmt.Errorf("Invalid MangoPlus client options: %s", err.Error()))
	}

	// The manga_viewer_v3 endpoint sets a cookie that is needed to get subscription-only images.
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	plus := &PlusClient{client: &client, options: options, baseAPI: BaseAPI, secret: options.Secret}

	plus.common.client = plus
	// Reuse the common client for the other services
	plus.Manga = (*MangaService)(&plus.common)
	plus.Page = (*PageService)(&plus.common)
	plus.User = (*UserService)(&plus.common)

	return plus
}

// Request: Sends a request to the MangaPlus API and decodes the protobuf response.
//
// path is relative to BaseAPI, all params are sent as query params as the API doesn't use bodies.
func (c *PlusClient) Request(ctx context.Context, method string, path string, params map[string]string) (*SuccessResult, error) {
	if err := c.ensureRegistered(ctx); err != nil {
		return nil, err
	}
	return c.request(ctx, method, path, params)
}

func (c *PlusClient) request(ctx context.Context, method string, path string, params map[string]string) (*SuccessResult, error) {
	u, err := url.JoinPath(c.baseAPI, path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = c.getFinalParams(params).Encode()
	req.Header.Set("User-Agent", c.options.UserAgent)

	if err := c.waitForNextRequest(ctx); err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MangaPlus API status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res Response
	if err := proto.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to decode MangaPlus response: %s", err.Error())
	}
	if res.Error != nil {
		return nil, fmt.Errorf("MangaPlus API error: %s", res.Error.message())
	}
	if res.Success == nil {
		return nil, fmt.Errorf("MangaPlus API unexpected error (Error and Success fields empty)")
	}
	return res.Success, nil
}

// waitForNextRequest: Sleeps until Options.RequestInterval has passed since the last request,
// or until ctx is done.
func (c *PlusClient) waitForNextRequest(ctx context.Context) error {
	c.lastRequestMu.Lock()
	defer c.lastRequestMu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Until(c.lastRequest.Add(c.options.RequestInterval))):
	}
	c.lastRequest = time.Now()
	return nil
}

// build and get the final params used for the request
func (c *PlusClient) getFinalParams(params map[string]string) url.Values {
	p := url.Values{}
	p.Set("os", "android")
	p.Set("os_ver", c.options.OSVersion)
	p.Set("app_ver", c.options.AppVersion)
	if c.secret != "" {
		p.Set("secret", c.secret)
	}
	for k, v := range params {
		p.Set(k, v)
	}
	return p
}

// ImageHeaders: Headers needed to download page images.
//
// Includes the cookie set by the last manga_viewer_v3 request, if any.
func (c *PlusClient) ImageHeaders(imageURL string) map[string]string {
	headers := map[string]string{"User-Agent": c.options.UserAgent}

	u, err := url.Parse(imageURL)
	if err != nil {
		return headers
	}
	var cookies []string
	for _, cookie := range c.client.Jar.Cookies(u) {
		cookies = append(cookies, cookie.String())
	}
	if len(cookies) != 0 {
		headers["Cookie"] = strings.Join(cookies, "; ")
	}
	return headers
}
