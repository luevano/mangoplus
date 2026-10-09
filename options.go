package mangoplus

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// The API rejects requests that don't look like they come from the Android app.
const defaultUserAgent = "okhttp/4.12.0"

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

	// Secret: Device secret from the official app, gives access to the account's subscription.
	//
	// Optional, if empty a new device is registered (free chapters only).
	Secret string
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
	if o.Viewer != ViewerLegacy && o.Viewer != ViewerApp {
		return fmt.Errorf("Viewer %q is not one of %q or %q", o.Viewer, ViewerLegacy, ViewerApp)
	}
	return nil
}

func DefaultOptions() Options {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic("Error while generating a random AndroidID.")
	}
	android_id := hex.EncodeToString(b)
	return Options{
		UserAgent:  defaultUserAgent,
		OSVersion:  "36",
		AppVersion: "265",
		AndroidID:  android_id,
		Viewer:     ViewerApp,
	}
}
