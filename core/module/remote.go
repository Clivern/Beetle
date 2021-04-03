// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ReleaseURL remote release URL
const ReleaseURL = "https://api.github.com/repos/Clivern/Beetle/releases/latest"

// LatestRelease struct
type LatestRelease struct {
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
}

// LoadFromJSON update object from json
func (lr *LatestRelease) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &lr)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (lr *LatestRelease) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&lr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetLatestRelease gets the latest beetle release
func GetLatestRelease() (LatestRelease, error) {
	result := LatestRelease{}

	httpClient := NewHTTPClient(20)

	response, err := httpClient.Get(
		context.TODO(),
		ReleaseURL,
		map[string]string{},
		map[string]string{},
	)

	if http.StatusOK != httpClient.GetStatusCode(response) || err != nil {
		return result, fmt.Errorf("Error: Unable to fetch latest release")
	}

	body, err := httpClient.ToString(response)

	if err != nil {
		return result, fmt.Errorf("Error: Unable to fetch latest release")
	}

	ok, err := result.LoadFromJSON([]byte(body))

	if !ok || err != nil {
		return result, fmt.Errorf("Error: Invalid remote response")
	}

	return result, nil
}
