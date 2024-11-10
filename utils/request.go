package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type APIResponse struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func DoRequest(method, url, payload, apiToken string) (*APIResponse, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var apiRes APIResponse
	if err := json.Unmarshal(body, &apiRes); err != nil {
		return nil, err
	}

	return &apiRes, nil
}
