package harvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type api struct {
	AccountID   string
	AccessToken string
	client      *http.Client
}

func NewClient(accountID string, accessToken string, client *http.Client) *api {
	if client == nil {
		client = &http.Client{}
	}

	return &api{
		AccountID:   accountID,
		AccessToken: accessToken,
		client:      client,
	}
}

func (a api) doRequest(method string, path string, reqObject any, respObject any) error {
	var body io.Reader
	if method == "PUT" || method == "POST" {
		jsonBytes, err := json.Marshal(reqObject)
		if err != nil {
			return fmt.Errorf("doRequest.Marshal:%w", err)
		}
		body = bytes.NewBuffer(jsonBytes)
	}

	url := "https://api.harvestapp.com" + path
	if strings.HasPrefix(path, "https://api.harvestapp.com") {
		url = path
	}

	req, err := http.NewRequest(method, url, body)
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

func (a api) GetTimeEntries() ([]TimeEntry, error) {
	var entries []TimeEntry
	var response TimeEntryResponse
	var err error

	for err = a.doRequest("GET", "/v2/time_entries", nil, response); err != nil && response.NextPage != ""; err = a.doRequest("GET", response.NextPage, nil, response) {
		for _, e := range response.TimeEntries {
			entries = append(entries, e)
		}
	}

	return entries, err
}
