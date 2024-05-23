package mangoplus

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"

type Options struct {
	UserAgent  string
	OSVersion  string
	AppVersion string
	AndroidID  string
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
		OSVersion:  "30",
		AppVersion: "133",
		AndroidID:  android_id,
	}
}
