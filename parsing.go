package rss

import (
    "encoding/xml"
    "time"
    "fmt"
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
    e.EncodeToken(xml.CharData(date.UTC().Format("Mon, 02 Jan 2006 15:04:05") + " GMT"))
    e.EncodeToken(xml.EndElement{start.Name})

    return
}

func (date *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var value string
    if err := d.DecodeElement(&value, &start); err != nil {
        return err
    }

    for _, tz := range([]string{"MST", "-0700"}) {
        for _, year := range([]string{"2006", "06"}) {
            for _, day := range([]string{"02", "2"}) {
                for _, dayOfWeek := range([]string{"Mon, ", ""}) {
                    format := fmt.Sprintf("%s%s Jan %s 15:04:05 %s", dayOfWeek, day, year, tz)
                    if date.tryParse(format, value) {
                        return nil
                    }
                }
            }
        }
    }

    for _, format := range([]string{
        "2006-01-02 15:04:05 -0700",
        "2006-01-02T15:04:05-07:00",
        "2006-01-02T15:04:05.000-07:00",
    }) {
        if date.tryParse(format, value) {
            return nil
        }
    }

    return fmt.Errorf("Can't parse date: %s.", value)
}

func (date *Date) tryParse(format string, value string) bool {
    var err error
    date.Time, err = time.Parse(format, value)
    return err == nil
}