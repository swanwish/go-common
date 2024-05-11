package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/swanwish/go-common/logs"
)

const (
	EmptyUsername = ""
	EmptyPassword = ""
)

func GetFullUrl(path, baseUrl string) string {
	if strings.Index(path, "://") != -1 {
		return path
	}
	protocolIndex := strings.Index(baseUrl, "://")
	if protocolIndex == -1 {
		logs.Errorf("The base url has no ://")
		return path
	}
	if path[0:1] == "/" {
		firstSlashIndex := strings.Index(baseUrl[protocolIndex+3:], "/")
		if firstSlashIndex != -1 {
			fullUrl := fmt.Sprintf("%s%s", baseUrl[:protocolIndex+3+firstSlashIndex], path)
			return fullUrl
		}
	} else {
		lastSlashIndex := strings.LastIndex(baseUrl, "/")
		if lastSlashIndex > protocolIndex+3 {
			fullUrl := fmt.Sprintf("%s%s", baseUrl[:lastSlashIndex+1], path)
			return fullUrl
		}
	}
	return path
}

func GetUrlContent(url string) (int, []byte, error) {
	logs.Debugf("Get content from url %s", url)
	response, err := http.Get(url)
	if err != nil {
		logs.Errorf("Failed to get content from url %s, the error is %v", url, err)
		return http.StatusInternalServerError, nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		logs.Errorf("Failed to read content from response, the error is %v\n", err)
		return http.StatusInternalServerError, nil, err
	}
	return response.StatusCode, contents, nil
}

func PostUrlContent(url string, content []byte, headers http.Header) (int, []byte, error) {
	return PostUrlContentWithBasicAuth(url, EmptyUsername, EmptyPassword, content, headers)
}

func PostUrlContentWithBasicAuth(url, username, password string, content []byte, headers http.Header) (int, []byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	if err != nil {
		logs.Errorf("Failed to create request, the error is %#v", err)
		return http.StatusBadRequest, nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if headers != nil {
		for key, value := range headers {
			req.Header[key] = value
		}
	}
	if username != "" || password != "" {
		req.SetBasicAuth(username, password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.Errorf("Failed to post content, the error is %v", err)
		return http.StatusInternalServerError, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, _ := io.ReadAll(resp.Body)
	logs.Debugf("The content is %s", string(body))
	return resp.StatusCode, body, nil
}

func PostRequest(url string, data url.Values, querys url.Values, headers http.Header) (int, []byte, error) {
	if headers != nil && headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	}
	return Request("POST", url, EmptyUsername, EmptyPassword, data.Encode(), querys.Encode(), headers)
}

func GetRequest(url string, data url.Values, querys url.Values, headers http.Header) (int, []byte, error) {
	if headers != nil && headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/json; charset=utf-8")
	}
	return Request("GET", url, EmptyUsername, EmptyPassword, data.Encode(), querys.Encode(), headers)
}

func PutRequest(url string, data url.Values, querys url.Values, headers http.Header) (int, []byte, error) {
	if headers != nil && headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/json; charset=utf-8")
	}
	return Request("PUT", url, EmptyUsername, EmptyPassword, data.Encode(), querys.Encode(), headers)
}

func Request(method, url, username, password, content, rawQuery string, headers http.Header) (int, []byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBufferString(content))
	if err != nil {
		logs.Errorf("Failed to get request, the error is %v", err)
		return http.StatusInternalServerError, nil, err
	}
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Add("Content-Length", strconv.Itoa(len(content)))
	if headers != nil {
		for key, value := range headers {
			req.Header[key] = value
		}
	}

	if rawQuery != "" {
		req.URL.RawQuery = rawQuery
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	resp, err := client.Do(req)
	if err != nil {
		logs.Errorf("Failed to post content, the error is %v", err)
		return http.StatusInternalServerError, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, _ := io.ReadAll(resp.Body)
	logs.Debugf("The content is %s", string(body))
	return resp.StatusCode, body, nil
}
