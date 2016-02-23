package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/swanwish/go-common/logs"
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

func GetUrlContent(url string) ([]byte, error) {
	logs.Debugf("Get content from url %s", url)
	response, err := http.Get(url)
	if err != nil {
		logs.Errorf("Failed to get content from url %s, the error is %v", url, err)
		return nil, err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logs.Errorf("Failed to read content from response, the error is %v\n", err)
		return nil, err
	}
	return contents, nil
}

func PostUrlContent(url string, content []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.Errorf("Failed to post content, the error is %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	logs.Debugf("The content is %s", string(body))
	return body, nil
}

func PostRequest(url string, data url.Values, headers http.Header) ([]byte, error) {
	return request("POST", url, data, headers)
}

func GetRequest(url string, data url.Values, headers http.Header) ([]byte, error) {
	return request("GET", url, data, headers)
}

func request(method, url string, data url.Values, headers http.Header) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		logs.Errorf("Failed to get request, the error is %v", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if headers != nil {
		for key, value := range headers {
			req.Header[key] = value
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		logs.Errorf("Failed to post content, the error is %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	logs.Debugf("The content is %s", string(body))
	return body, nil
}
