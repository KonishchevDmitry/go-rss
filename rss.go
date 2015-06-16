package rss

import (
    "bytes"
    "encoding/xml"
    "errors"
    "fmt"
    "io"
    "mime"
    "net/http"
    "time"

    "code.google.com/p/go-charset/charset"
    _ "code.google.com/p/go-charset/data"
)

type GetParams struct {
    Timeout time.Duration
}

var DefaultGetParams = GetParams{
    Timeout: 30 * time.Second,
}

const ContentType = "application/rss+xml"

func Get(url string) (*Feed, error) {
    return GetWithParams(url, DefaultGetParams)
}

func GetWithParams(url string, params GetParams) (*Feed, error) {
    client := &http.Client{
        Timeout: params.Timeout,
    }

    response, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    if err := checkResponse(response); err != nil {
        return nil, err
    }

    return Read(response.Body)
}

func Read(reader io.Reader) (*Feed, error) {
    rss := rssRoot{}

    decoder := xml.NewDecoder(reader)
    decoder.CharsetReader = charset.NewReader
    if err := decoder.Decode(&rss); err != nil {
        return nil, err
    }

    switch rss.Version {
        case "2.0", "0.92", "0.91":
        default:
            return nil, fmt.Errorf("Invalid RSS version: %s.", rss.Version)
    }

    if rss.Channel == nil {
        return nil, fmt.Errorf("The document doesn't conform to RSS specification.")
    }

    return rss.Channel, nil
}

func Parse(data []byte) (*Feed, error) {
    return Read(bytes.NewReader(data))
}

func Write(feed *Feed, writer io.Writer) error {
    if _, err := writer.Write([]byte(xml.Header)); err != nil {
        return err
    }

    rss := rssRoot{Version: "2.0", Channel: feed}
    encoder := xml.NewEncoder(writer)
    encoder.Indent("", "    ")
    return encoder.Encode(&rss)
}

func Generate(feed *Feed) ([]byte, error) {
    var buffer bytes.Buffer

    if err := Write(feed, &buffer); err != nil {
        return nil, err
    }

    return buffer.Bytes(), nil
}

func checkResponse(response *http.Response) error {
    if response.StatusCode != http.StatusOK {
        return errors.New(response.Status)
    }

    contentType := response.Header.Get("Content-Type")
    mediaType, _, err := mime.ParseMediaType(contentType)
    if err != nil {
        return fmt.Errorf("The feed has an invalid Content-Type: %s", err)
    }

    allowedMediaTypes := map[string]bool {ContentType: true, "application/xml": true, "text/xml": true}
    if !allowedMediaTypes[mediaType] {
        return fmt.Errorf("The feed has an invalid Content-Type (%s).", mediaType)
    }

    return nil
}