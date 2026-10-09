# mangoplus

Golang API wrapper for MangaPlus API and MangaPlusCreators API.

There's no API documentation so most of the implementation is just winging it.

> **Warning**
> 
> The API implementation is not stable and may change at any time.

## Installation

To install, do `go get -u github.com/luevano/mangoplus@latest`.

## Usage

Basic usage example.

```go
package main

import (
    "fmt"

    "github.com/luevano/mangoplus"
)

func main() {
    // Create new client, a device is registered on the first request.
    c, err := mangoplus.NewPlusClient(mangoplus.DefaultOptions())
    if err != nil {
        panic(err)
    }

    // Get manga by id (Ghost Fixers).
    mangaDetails, err := c.Manga.Get("100310")
    if err != nil {
        panic(err)
    }
    for _, chapter := range mangaDetails.ChapterListV2 {
        fmt.Println(chapter.Name)
    }
}
```

Errors can be checked with `errors.Is` against `ErrNotFound`, `ErrRateLimited`, `ErrMaintenance`,
`ErrSubscriptionRequired` and `ErrRegionBlocked`, or inspected with `errors.As` as `*HTTPError` (non-200
status) or `*ServerError` (error result from the API, with its `Code`).

### Free device

Without a secret, a new free device is registered on the first request. To reuse it in later runs
instead of registering a new one every time, save the secret passed to `Options.OnRegister` and give it
back as `Options.FreeSecret`.

### Subscription

To read subscription chapters, set `Options.Secret` to the device secret of the official app where the
subscription is active. See [docs/subscription.md](docs/subscription.md) for how to get it.

`ViewerApp` (the default) is the only viewer that sets the cookie subscription-only images need.
`ViewerLegacy` uses the old endpoint and is deprecated.

## Known issues

- **MangaPlus Creators** (`creators` package) doesn't work: its API answers every request with an AWS
  WAF bot challenge (`202` with `x-amzn-waf-action: challenge`). Checked 2026-10-09.

## Development

The API uses protobuf, the schema is in `mangaplus.proto`. After changing it, run `go generate`, which
generates `mangaplus.pb.go` with [buf](https://buf.build) and
[protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) (both through `go run`).

Tests run offline against responses saved in `testdata/`. To run the tests against the live API, set
`MANGOPLUS_LIVE=1`.

## Credits

- [frank_mangaplus](https://github.com/akitaonrails/frank_mangaplus) by Fabio Akita (MIT): the
  protobuf field numbers, the device registration, the viewer params and cookie, and the guide adapted in
  [docs/subscription.md](docs/subscription.md).
- [hyugogirubato/mangaplus](https://github.com/hyugogirubato/mangaplus) (GPL-3.0): used only as a
  reference to cross-check endpoints and field numbers, no code taken.

## Contributing

Any contributions are welcome.
