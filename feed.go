package rss

import (
    "fmt"
    "time"
)

type Feed struct {
    Title       string   `xml:"title"`
    Link        string   `xml:"link"`
    Description string   `xml:"description"`
    Image       *Image   `xml:"image"`
    Language    string   `xml:"language,omitempty"`
    Date        Date     `xml:"pubDate"`
    Category    []string `xml:"category"`
    Generator   string   `xml:"generator,omitempty"`
    Ttl         int      `xml:"ttl,omitempty"`
    Items       []*Item  `xml:"item"`
}

type Image struct {
    Url    string `xml:"url"`
    Title  string `xml:"title"`
    Link   string `xml:"link"`
    Width  int    `xml:"width,omitempty"`
    Height int    `xml:"height,omitempty"`
}

type Date struct {
    time.Time
}

type Item struct {
    Title       string       `xml:"title,omitempty"`
    Guid        Guid         `xml:"guid"`
    Link        string       `xml:"link,omitempty"`
    Description string       `xml:"description,omitempty"`
    Enclosure   []*Enclosure `xml:"enclosure"`
    Comments    string       `xml:"comments,omitempty"`
    Date        Date         `xml:"pubDate"`
    Author      string       `xml:"author,omitempty"`
    Category    []string     `xml:"category"`
}

type Guid struct {
    Id          string `xml:",chardata"`
    IsPermaLink *bool  `xml:"isPermaLink,attr,omitempty"`
}

type Enclosure struct {
    Url    string `xml:"url,attr"`
    Type   string `xml:"type,attr"`
    Length int    `xml:"length,attr"`
}

func (feed *Feed) String() string {
    if feed == nil {
        return fmt.Sprintf("%#v", feed)
    }

    xml, err := Generate(feed)
    if err == nil {
        return string(xml)
    }

    return fmt.Sprintf("XML generation error: %s. Go representation: %#v", err, feed)
}