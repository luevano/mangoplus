# mangoplus

Golang API wrapper for MangaPlus API.

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
    "net/url"

    "github.com/luevano/mangoplus"
)

func main() {
    // Create new client.
    c := mangoplus.NewPlusClient()

    // ID for Ghost Fixers.
    id := 100310

    // Get manga by id.
    mangaDetails, err := c.Manga.Get(id)
    if err != nil {
        panic(err)
    }
    for _, chapter := range mangaDetails.ChapterListGroup {
        for _, c := range chapter.FirstChapterList {
            fmt.Println(c.Name)
        }
    }
}
```

## Contributing

Any contributions are welcome.
