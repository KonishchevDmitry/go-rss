package rss

import (
    "bytes"
    "encoding/xml"
    "errors"
    "fmt"
    "io"
    "net/http"
    "time"
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

    if response.StatusCode != http.StatusOK {
        return nil, errors.New(response.Status)
    }

    content_type := response.Header.Get("Content-Type")
    if content_type != ContentType {
        return nil, errors.New("The feed has an invalid Content-Type.")
    }

    return Read(response.Body)
}

func Read(reader io.Reader) (*Feed, error) {
    rss := rssRoot{}
    if err := xml.NewDecoder(reader).Decode(&rss); err != nil {
        return nil, err
    }

    if rss.Version != "2.0" {
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