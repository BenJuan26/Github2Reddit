package github

import (
	"encoding/json"
	"net/http"
)

type Release struct {
	ID         string `json:"id"`
	TagName    string `json:"tag_name"`
	Name       string `json:"name,omitempty"`
	Prerelease bool   `json:"prerelease"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
}

func ParseReleasePayload(req *http.Request) Release {
	type releaseBody struct {
		Action  string  `json:"action"`
		Release Release `json:"release"`
	}
	var body releaseBody
	json.NewDecoder(req.Body).Decode(&body)

	release := body.Release
	return release
}
