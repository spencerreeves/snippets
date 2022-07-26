package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type api struct {
	AccountID   string
	AccessToken string
	client      *http.Client
	URL         string
}

func NewClient(accountID string, accessToken string, client *http.Client, url string) *api {
	if client == nil {
		client = &http.Client{}
	}
	if url == "" {
		url = "https://api.harvestapp.com/"
	}

	return &api{
		AccountID:   accountID,
		AccessToken: accessToken,
		client:      client,
		URL:         url,
	}
}

func (a api) DoRequest(method string, path string, reqObject any, respObject any) error {
	var body io.Reader
	if method == "PUT" || method == "POST" {
		jsonBytes, err := json.Marshal(reqObject)
		if err != nil {
			return fmt.Errorf("doRequest.Marshal:%w", err)
		}
		body = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, a.URL+path, body)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	req.Header.Set("Harvest-Account-ID", a.AccountID)
	req.Header.Set("Authorization", "Bearer "+a.AccessToken)

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("non-ok response: %s", respBody)
	}

	if respObject != nil {
		err = json.Unmarshal(respBody, respObject)
		if err != nil {
			return fmt.Errorf("unable to unmarshal response: %w", err)
		}
	}

	return nil
}
