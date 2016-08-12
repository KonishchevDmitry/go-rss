package rss

import (
    "bytes"
    "crypto/tls"
    "encoding/xml"
    "errors"
    "fmt"
    "io"
    "mime"
    "net/http"
    "strings"
    "time"

    "golang.org/x/net/html/charset"
)

type GetParams struct {
    Timeout time.Duration
    Cookies []*http.Cookie
    SkipContentTypeCheck bool
    SkipCertificateCheck bool
}

const ContentType = "application/rss+xml"

func Get(url string) (*Feed, error) {
    return GetWithParams(url, GetParams{})
}

func GetWithParams(url string, params GetParams) (feed *Feed, err error) {
    client := ClientFromParams(params)
    allowedMediaTypes := []string{"application/rss+xml", "application/xml", "text/xml"}

    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return
    }

    // Some servers don't return the feed without Accept header. For example one server (which requires cookie
    // authentication) always returns login page if Accept header is not specified.
    request.Header.Set("Accept", strings.Join(allowedMediaTypes, ", "))

    for _, cookie := range params.Cookies {
        request.AddCookie(cookie)
    }

    response, err := client.Do(request)
    if err != nil {
        return
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return nil, errors.New(response.Status)
    }

    if !params.SkipContentTypeCheck {
        err = checkContentType(response, allowedMediaTypes)
        if err != nil {
            return
        }
    }

    return Read(response.Body)
}

func ClientFromParams(params GetParams) (*http.Client) {
    client := &http.Client{}

    if params.Timeout == 0 {
        client.Timeout = 30 * time.Second
    } else {
        client.Timeout = params.Timeout
    }

    if params.SkipCertificateCheck {
        client.Transport = &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
        }

    }

    return client
}

func Read(reader io.Reader) (*Feed, error) {
    rss := rssRoot{}

    decoder := xml.NewDecoder(reader)
    decoder.CharsetReader = charset.NewReaderLabel
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

func checkContentType(response *http.Response, allowedMediaTypes []string) error {
    contentType := response.Header.Get("Content-Type")
    mediaType, _, err := mime.ParseMediaType(contentType)
    if err != nil {
        return fmt.Errorf("The feed has an invalid Content-Type: %s", err)
    }

    for _, allowedMediaType := range allowedMediaTypes {
        if mediaType == allowedMediaType {
            return nil
        }
    }

    return fmt.Errorf("The feed has an invalid Content-Type (%s).", mediaType)
}