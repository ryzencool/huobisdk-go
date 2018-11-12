package robot

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36"

var tr = &http.Transport{
	//Proxy:           http.ProxyURL(httpProxy),
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var client = &http.Client{
	Transport: tr,
	Timeout:   time.Second * 3,
}

// Get is http util to send get request
func Get(webURL string, params map[string]string) (string, error) {
	if params != nil || len(params) > 0 {
		var paramArr []string
		for k, v := range params {
			paramArr = append(paramArr, k+"="+v)
		}
		webURL = webURL + "?" + strings.Join(paramArr, "&")
	}
	request, err := http.NewRequest(http.MethodGet, webURL, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "x-www-form-urlencoded")
	request.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(request)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}

// Post is send post request
func Post(url string, params map[string]interface{}) (string, error) {
	formed, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(formed))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(request)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
