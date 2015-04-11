package rss

import (
    "encoding/xml"
    "net/http"
)

type rssRoot struct {
    XMLName xml.Name `xml:"rss"`
    Version string   `xml:"version,attr"`
    Channel *Feed    `xml:"channel"`
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