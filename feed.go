package rss

import (
    "bytes"
    "encoding/xml"
    "errors"
    "fmt"
    "net/http"
    "time"
)

type rssRoot struct {
    XMLName xml.Name `xml:"rss"`
    Version string   `xml:"version,attr"`
    Channel *Feed    `xml:"channel"`
}

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


func Parse(data []byte) (*Feed, error) {
    rss := rssRoot{}

    if err := xml.Unmarshal(data, &rss); err != nil {
        return nil, err
    }

    if rss.Version != "2.0" {
        return nil, errors.New(fmt.Sprintf("Invalid RSS version: %s.", rss.Version))
    }

    if rss.Channel == nil {
        return nil, errors.New("The document doesn't conform to RSS specification.")
    }

    return rss.Channel, nil
}

func Generate(feed *Feed) ([]byte, error) {
    rss := rssRoot{Version: "2.0", Channel: feed}

    data, err := xml.MarshalIndent(&rss, "", "    ")
    if err != nil {
        return nil, err
    }

    rssData := bytes.NewBufferString(xml.Header)
    rssData.Write(data)

    return rssData.Bytes(), nil
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


func (guid *Guid) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
    if guid.Id == "" {
        return
    }

    if guid.IsPermaLink != nil {
        value := "true"
        if !*guid.IsPermaLink {
            value = "false"
        }

        attr := xml.Attr{
            Name: xml.Name{Local: "isPermaLink"},
            Value: value,
        }

        start.Attr = append(start.Attr, attr)
    }

    e.EncodeToken(start)
    e.EncodeToken(xml.CharData(guid.Id))
    e.EncodeToken(xml.EndElement{start.Name})

    return
}

func (date *Date) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
    if date.IsZero() {
        return
    }

    e.EncodeToken(start)
    e.EncodeToken(xml.CharData(date.Format(http.TimeFormat)))
    e.EncodeToken(xml.EndElement{start.Name})

    return
}

func (date *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var value string

    if err := d.DecodeElement(&value, &start); err != nil {
        return err
    }

    time, err := http.ParseTime(value)
    if err != nil {
        return err
    }

    date.Time = time
    return nil
}