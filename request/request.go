package request

import (
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var FakeHeaders = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Charset":  "UTF-8,*;q=0.5",
	"Accept-Encoding": "gzip,deflate,sdch",
	"Accept-Language": "en-US,en;q=0.8",
	"User-Agent":      "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0",
}

// Get Response/Cookie
func Request(method, url string, headers map[string]string) (*http.Response, map[string]string, error) {
	// Create New Request [HTTP]
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DisableCompression:  true,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, nil, err
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Minute,
		Jar:       jar,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, nil, err
	}

	for k, v := range FakeHeaders {
		req.Header.Set(k, v)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	for _, ok := headers["Referer"]; !ok; {
		req.Header.Set("Referer", url)
	}
	res, err := client.Do(req)

	cookies := map[string]string{}
	for _, cookie := range res.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	return res, cookies, err
}

// Get HTML DOC
func Get(url, refer string, headers map[string]string) (string, map[string]string, error) {
	body, res_headers, err := GetByte(url, refer, headers)
	return string(body), res_headers, err
}

// Get Byte Data * Video Data & Image Data
func GetByte(url, refer string, headers map[string]string) ([]byte, map[string]string, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	if refer != "" {
		headers["Referer"] = refer
	}
	res, res_headers, err := Request(http.MethodGet, url, headers)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close() // nolint

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
	case "deflate":
		reader = flate.NewReader(res.Body)
	default:
		reader = res.Body
	}
	defer reader.Close() // nolint

	body, err := io.ReadAll(reader)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}
	return body, res_headers, nil
}

// Get Headers Data
// * Content-Type
// * Content-Length
func Headers(url, refer string) (http.Header, map[string]string, error) {
	headers := map[string]string{
		"Referer": refer,
	}
	res, res_headers, err := Request(http.MethodGet, url, headers)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close() // nolint
	return res.Header, res_headers, nil
}
