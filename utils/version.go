package utils

import (
	"encoding/json"
	"errors"
	"html-cli/types"
	"io"

	"net/http"
)

func GetLatestRelease() (*types.GithubRelease, error) {
	resp, err := http.Get("https://api.github.com/repos/CyberL1/html-cli/releases/latest")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 403 {
		return nil, errors.New("rate limited by github")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	release := &types.GithubRelease{}
	err = json.Unmarshal(body, release)
	if err != nil {
		return nil, err
	}
	return release, nil
}
