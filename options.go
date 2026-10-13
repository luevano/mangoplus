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

type Options struct {
	UserAgent  string
	OSVersion  string
	AppVersion string
	AndroidID  string
	Viewer     Viewer

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
	return nil
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

		RequestInterval: 500 * time.Millisecond,
	}
}
