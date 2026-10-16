package mangoplus

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// The API rejects requests that don't look like they come from the Android app.
const defaultUserAgent = "okhttp/4.12.0"

// defaultWebUserAgent: The web api gets a browser's, a copy of libmangal's netutil.BrowserUserAgent
// like creators' default.
const defaultWebUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:158.0) Gecko/20100101 Firefox/158.0"

// defaultTimeout: Timeout of the client used when Options.HTTPClient is nil.
const defaultTimeout = time.Minute

// Viewer: Which endpoint and params are used to get chapter pages.
type Viewer string

const (
	// ViewerApp: Same endpoint and params as the official app.
	//
	// The only one that sets the cookie needed to download subscription-only images.
	ViewerApp Viewer = "app"
	// ViewerLegacy: Old endpoint, without the reading params the app sends.
	//
	// Deprecated: Use ViewerApp, kept in case the old endpoint behaves differently.
	ViewerLegacy Viewer = "legacy"
)

// API: Which MangaPlus api the client talks to.
type API string

const (
	// APIApp: The Android app's api (BaseAPI), with a device secret (registered if none is given).
	APIApp API = "app"
	// APIWeb: The website's api (BaseWebAPI), no secret or device: only what the website shows
	// (first and latest chapters of a title, no account), at most 1 request per second.
	APIWeb API = "web"
	// APIAuto: APIWeb when neither Secret nor FreeSecret is set, APIApp otherwise.
	APIAuto API = "auto"
)

type Options struct {
	UserAgent  string
	OSVersion  string
	AppVersion string
	AndroidID  string
	Viewer     Viewer

	// API: Optional, APIApp when empty.
	API API

	// WebUserAgent: User agent of the web api's requests and their images.
	//
	// Optional, a browser's when empty.
	WebUserAgent string

	// RequestInterval: Minimum time between API requests.
	//
	// Warning: lowering or removing it (0) risks a lockout, the API is reported to block the device
	// for 10+ minutes after ~10 requests in a few seconds.
	RequestInterval time.Duration

	// Secret: Device secret from the official app, gives access to the account's subscription.
	//
	// Optional, takes precedence over FreeSecret.
	Secret string

	// FreeSecret: Device secret of a previously registered free device (free chapters only).
	//
	// Optional, if both secrets are empty a new device is registered on the first request.
	FreeSecret string

	// OnRegister: Called with the secret of a newly registered device, e.g. to save it as FreeSecret.
	OnRegister func(secret string)

	// HTTPClient: Client used for every request, copied to give it a cookie jar if it has none.
	//
	// Optional, a client with a 1 minute timeout is used when nil.
	HTTPClient *http.Client
}

func (o Options) validate() error {
	if o.UserAgent == "" {
		return fmt.Errorf("UserAgent is empty")
	}
	if o.OSVersion == "" {
		return fmt.Errorf("OSVersion is empty")
	}
	if o.AppVersion == "" {
		return fmt.Errorf("AppVersion is empty")
	}
	if o.AndroidID == "" {
		return fmt.Errorf("AndroidID is empty")
	}
	if o.RequestInterval < 0 {
		return fmt.Errorf("RequestInterval is negative")
	}
	if o.Viewer != ViewerLegacy && o.Viewer != ViewerApp {
		return fmt.Errorf("Viewer %q is not one of %q or %q", o.Viewer, ViewerLegacy, ViewerApp)
	}
	switch o.API {
	case "", APIApp, APIWeb, APIAuto:
	default:
		return fmt.Errorf("API %q is not one of %q, %q or %q", o.API, APIApp, APIWeb, APIAuto)
	}
	return nil
}

// useWebAPI: Whether these options pick the web api, see API.
func (o Options) useWebAPI() bool {
	switch o.API {
	case APIWeb:
		return true
	case APIAuto:
		return o.Secret == "" && o.FreeSecret == ""
	}
	return false
}

func DefaultOptions() Options {
	b := make([]byte, 8)
	rand.Read(b)
	android_id := hex.EncodeToString(b)
	return Options{
		UserAgent:  defaultUserAgent,
		OSVersion:  "36",
		AppVersion: "265",
		AndroidID:  android_id,
		Viewer:     ViewerApp,
		API:        APIApp,

		RequestInterval: 500 * time.Millisecond,
	}
}
