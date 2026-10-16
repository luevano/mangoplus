// Package mangoplus provides an API wrapper for MangaPlus API.
package mangoplus

//go:generate go run github.com/bufbuild/buf/cmd/buf@v1.73.0 generate --template buf.gen.yaml

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	BaseAPI = "https://jumpg-api.tokyo-cdn.com/api"
	// BaseWebAPI is the web reader's api, limited to what the browser gets
	BaseWebAPI = "https://jumpg-webapi.tokyo-cdn.com/api"
)

// webRequestInterval: Minimum time between web api requests, from keiyoushi's MangaPlus
// extension (ref/keiyoushi-extensions src/all/mangaplus/.../MangaPlus.kt:
// rateLimit(1) { it.host == apiUrlHost }, API_URL being jumpg-webapi).
const webRequestInterval = time.Second

// webOrigin: The website, sent as Origin/Referer of web api requests like the browser does.
const webOrigin = "https://mangaplus.shueisha.co.jp"

// webContentLanguage: The web api's clang (language of its messages), as keiyoushi sends it.
const webContentLanguage = "eng"

// maxBodySize caps the responses read.
const maxBodySize = 100 << 20

// PlusClient: The MangaPlus client.
type PlusClient struct {
	client  *http.Client
	options Options
	// BaseAPI, only changed in tests
	baseAPI string

	// Device secret sent on every request, from Options.Secret, Options.FreeSecret or registering a new device.
	secret     string
	registerMu sync.Mutex

	// web api: the browser's random SESSION-TOKEN, and the last viewer's token for the images
	web          bool
	sessionToken string
	viewToken    string
	viewTokenMu  sync.Mutex

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
// Use DefaultOptions for defaults, all fields except the secrets and OnRegister must be non-empty.
//
// If no secret is provided, a new device is registered on the first request, unless
// Options.API picks the web api, which needs none.
func NewPlusClient(options Options) (*PlusClient, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("mangoplus: invalid options: %w", err)
	}

	// The manga_viewer_v3 endpoint sets a cookie that is needed to get subscription-only images.
	client := http.Client{Timeout: defaultTimeout}
	if options.HTTPClient != nil {
		client = *options.HTTPClient
	}
	if client.Jar == nil {
		client.Jar, _ = cookiejar.New(nil)
	}
	secret := options.Secret
	if secret == "" {
		secret = options.FreeSecret
	}
	plus := &PlusClient{client: &client, options: options, baseAPI: BaseAPI, secret: secret}
	if options.useWebAPI() {
		plus.web = true
		plus.baseAPI = BaseWebAPI
		plus.secret = ""
		plus.sessionToken = newSessionToken()
		if plus.options.WebUserAgent == "" {
			plus.options.WebUserAgent = defaultWebUserAgent
		}
		plus.options.RequestInterval = max(plus.options.RequestInterval, webRequestInterval)
	}

	plus.common.client = plus
	// Reuse the common client for the other services
	plus.Manga = (*MangaService)(&plus.common)
	plus.Page = (*PageService)(&plus.common)
	plus.User = (*UserService)(&plus.common)

	return plus, nil
}

// WebAPI: Whether the client uses the web api (BaseWebAPI), see Options.API.
func (c *PlusClient) WebAPI() bool {
	return c.web
}

// Request: Sends a request to the MangaPlus API and decodes the protobuf response.
//
// path is relative to BaseAPI (BaseWebAPI with the web api), all params are sent as query params
// as the API doesn't use bodies.
func (c *PlusClient) Request(ctx context.Context, method string, path string, params map[string]string) (*SuccessResult, error) {
	if !c.web {
		if err := c.ensureRegistered(ctx); err != nil {
			return nil, err
		}
	}
	return c.request(ctx, method, path, params)
}

func (c *PlusClient) request(ctx context.Context, method string, path string, params map[string]string) (*SuccessResult, error) {
	// prefix for errors that aren't from the API itself
	op := fmt.Sprintf("mangoplus: %s /%s", method, path)

	u, err := url.JoinPath(c.baseAPI, path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	req.URL.RawQuery = c.getFinalParams(params).Encode()
	if c.web {
		req.Header.Set("User-Agent", c.options.WebUserAgent)
		req.Header.Set("Session-Token", c.sessionToken)
		req.Header.Set("Origin", webOrigin)
		req.Header.Set("Referer", webOrigin+"/")
	} else {
		req.Header.Set("User-Agent", c.options.UserAgent)
	}

	if err := c.waitForNextRequest(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		// net/http puts the full url in its errors, query included
		var urlErr *url.Error
		if errors.As(err, &urlErr) && os.Getenv(UnsafeShowSecretEnv) != "1" {
			urlErr.URL = redactURL(urlErr.URL)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{Method: method, Path: "/" + path, StatusCode: resp.StatusCode}
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize+1))
	if err == nil && len(body) > maxBodySize {
		err = fmt.Errorf("response over %d bytes", maxBodySize)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: read response: %w", op, err)
	}
	var res Response
	if err := proto.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("%s: decode response: %w", op, err)
	}
	if res.Error != nil {
		return nil, newServerError(method, "/"+path, res.Error)
	}
	if res.Success == nil {
		return nil, fmt.Errorf("%s: empty response", op)
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
	if !c.web {
		p.Set("os", "android")
		p.Set("os_ver", c.options.OSVersion)
		p.Set("app_ver", c.options.AppVersion)
	}
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
// Includes the cookie set by the last manga_viewer_v3 request, if any, or with the web api the
// last viewer's token.
func (c *PlusClient) ImageHeaders(imageURL string) map[string]string {
	headers := map[string]string{"User-Agent": c.options.UserAgent}
	if c.web {
		headers["User-Agent"] = c.options.WebUserAgent
		headers["Referer"] = webOrigin + "/"
		c.viewTokenMu.Lock()
		if c.viewToken != "" {
			headers["Plus-Vw-Token"] = c.viewToken
		}
		c.viewTokenMu.Unlock()
	}

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

// newSessionToken: A random uuid (v4), what the website sends as SESSION-TOKEN.
func newSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	h := hex.EncodeToString(b)
	return h[:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
}

// UnsafeShowSecretEnv: Set it to "1" to show the secret and registration keys in errors instead of
// "REDACTED". Unsafe, anyone who sees the error can use the account.
const UnsafeShowSecretEnv = "MANGOPLUS_UNSAFE_SHOW_SECRET"

// params that identify the device and give access to the account
var sensitiveParams = []string{"secret", "device_token", "security_key"}

// redactURL: Hides the sensitive params of the url so it can be shown in errors.
func redactURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		// can't tell what's in it, drop the query altogether
		before, _, _ := strings.Cut(rawURL, "?")
		return before
	}
	return redact(u)
}

// RedactURL: u with the secret and registration keys hidden, for logs (e.g. a debug
// transport, see libmangal's netutil.LoggingTransport), unless UnsafeShowSecretEnv is "1".
func RedactURL(u *url.URL) string {
	if os.Getenv(UnsafeShowSecretEnv) == "1" {
		return u.String()
	}
	return redact(u)
}

func redact(u *url.URL) string {
	redacted := *u
	query := u.Query()
	for _, param := range sensitiveParams {
		if query.Has(param) {
			query.Set(param, "REDACTED")
		}
	}
	redacted.RawQuery = query.Encode()
	return redacted.String()
}
