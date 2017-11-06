package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type SubmitBody struct {
	Kind      string `json:"kind"`
	Subreddit string `json:"sr"`
	Text      string `json:"text,omitempty"`
	Title     string `json:"title"`
	URL       string `json:"url,omitempty"`
}

func GetToken(user, pass, clientID, clientSecret, botName string) (string, error) {
    tokenReqBaseURL := "https://www.reddit.com/api/v1/access_token"
    tokenReqURL := fmt.Sprintf("%s?grant_type=password&username=%s&password=%s", tokenReqBaseURL, user, pass)

    tokenReq, err := http.NewRequest("POST", tokenReqURL, nil)
    tokenReq.Header.Add("User-Agent", botName)
    tokenReq.SetBasicAuth(clientID, clientSecret)

    client := &http.Client{}
    resp, err := client.Do(tokenReq)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return "", errors.New("Unexpected status code " + resp.Status)
    }

	var tokenResp TokenResponse
	json.NewDecoder(resp.Body).Decode(&tokenResp)

	return tokenResp.AccessToken, nil
}
