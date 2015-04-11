Simple RSS parser and generator for Go

Usage example:
```go
package main

import (
    "fmt"
    "github.com/KonishchevDmitry/go-rss"
)

var rssData = `
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
    <channel>
        <title>Feed title</title>
        <link>http://example.com/</link>
        <description>Feed description</description>
        <item>
            <title>Item 1</title>
            <link>http://example.com/item1</link>
            <description>Item 1 description</description>
        </item>
    </channel>
</rss>
`

func main() {
    url := "http://example.com/feed.rss"

    // Get an RSS feed from URL
    feed, err := rss.Get(url)
    if err != nil {
        fmt.Printf("Failed to get %s: %s\n", url, err)
    }

    // Get an RSS feed from string
    feed, err = rss.Parse([]byte(rssData))
    if err != nil {
        fmt.Println("Parsing error:", err)
        return
    }

    // Change <generator> element value
    feed.Generator = "go-rss"

    // Generate the modified RSS feed
    data, err := rss.Generate(feed)
    if err != nil {
        fmt.Println("RSS generation error:", err)
        return
    }

    fmt.Printf("%s\n", data)
}
```

Output:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
    <channel>
        <title>Feed title</title>
        <link>http://example.com/</link>
        <description>Feed description</description>
        <generator>go-rss</generator>
        <item>
            <title>Item 1</title>
            <link>http://example.com/item1</link>
            <description>Item 1 description</description>
        </item>
    </channel>
</rss>
```