package httpRequest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"network-detector/libs/logger"
	"time"
)

func Get(url string) (result map[string]interface{}, err error) {
	logger.Printf("http request get url: %s\n", url)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		logger.Printf("http request get error: %s\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("http request read body error: %s\n", err.Error())
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Printf("json unmarshal error: %s\n", err.Error())
		return nil, err
	}

	return result, nil
}

// contentType: application/json
func Post(url string, data interface{}, contentType string) (content string, err error) {
	logger.Printf("http request post url: %s, data: %v\n", url, data)
	jsonStr, err := json.Marshal(data)
	if err != nil {
		logger.Printf("http request marshal JSON error: %s\n", err.Error())
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Printf("http request create request error: %s\n", err.Error())
		return "", err
	}
	req.Header.Add("content-type", contentType)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Printf("http request post error: %s\n", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("http request read body error: %s\n", err.Error())
		return "", err
	}

	content = string(body)
	return content, nil
}
