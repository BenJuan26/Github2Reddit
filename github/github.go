package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Release struct {
	URL        string `json:"html_url"`
	ID         int    `json:"id"`
	TagName    string `json:"tag_name"`
	Name       string `json:"name,omitempty"`
	Prerelease bool   `json:"prerelease"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
}

type Repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	URL      string `json:"url"`
}

type ReleaseRequestBody struct {
	Action     string     `json:"action"`
	Release    Release    `json:"release"`
	Repository Repository `json:"repository"`
}


func ParseReleasePayload(req *http.Request) (ReleaseRequestBody, error) {
	var body ReleaseRequestBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		fmt.Println("decoder error")
		return ReleaseRequestBody{}, err
	}

	return body, nil
}
