package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type MoodleAPI struct {
	baseURL  string
	endpoint string
}

func NewMoodleAPI() *MoodleAPI {
	return &MoodleAPI{
		baseURL:  "https://moodle.astanait.edu.kz",
		endpoint: "/webservice/rest/server.php",
	}
}

func (self *MoodleAPI) MakeRequest(function string, token string, params map[string]string) (map[string]any, error) {
	urlParams := url.Values{
		"moodlewsrestformat": {"json"},
		"wstoken":            {token},
		"wsfunction":         {function},
	}

	for key, value := range params {
		urlParams.Add(key, value)
	}

	reqURL := fmt.Sprintf("%s%s?%s", self.baseURL, self.endpoint, urlParams.Encode())

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
