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
    e.EncodeToken(xml.CharData(date.Format("Mon, 02 Jan 2006 15:04:05 GMT")))
    e.EncodeToken(xml.EndElement{start.Name})

    return
}

func (date *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var value string
    if err := d.DecodeElement(&value, &start); err != nil {
        return err
    }

    for _, year := range([]string{"2006", "06"}) {
        for _, day := range([]string{"02", "2"}) {
            for _, dayOfWeek := range([]string{"Mon, ", ""}) {
                dateFormat := fmt.Sprintf("%s%s Jan %s 15:04:05 MST", dayOfWeek, day, year)
                dateTime, err := time.Parse(dateFormat, value)
                if err == nil {
                    date.Time = dateTime
                    return nil
                }
            }
        }
    }

    return fmt.Errorf("Can't parse date: %s.", value)
}