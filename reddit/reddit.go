package reddit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	//"strings"
	"text/template"

	"github.com/BenJuan26/Github2Reddit/github"
	"github.com/gorilla/schema"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type SubmitBody struct {
	Kind      string `json:"kind" schema:"kind"`
	Subreddit string `json:"subreddit" schema:"sr"`
	Text      string `json:"text,omitempty" schema:"text,omitempty"`
	Title     string `json:"title" schema:"title"`
	URL       string `json:"url,omitempty" schema:"url,omitempty"`
	APIType   string `json:"api_type,omitempty" schema:"api_type,omitempty"`
}

func GetToken(user, pass, clientID, clientSecret, botName string) (Token, error) {
    tokenReqBaseURL := "https://www.reddit.com/api/v1/access_token"
    tokenReqURL := fmt.Sprintf("%s?grant_type=password&username=%s&password=%s", tokenReqBaseURL, user, pass)

    tokenReq, err := http.NewRequest("POST", tokenReqURL, nil)
    tokenReq.Header.Add("User-Agent", botName)
    tokenReq.SetBasicAuth(clientID, clientSecret)

    client := &http.Client{}
    resp, err := client.Do(tokenReq)
    if err != nil {
        return Token{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return Token{}, errors.New("Unexpected status code " + resp.Status)
    }

	var token Token
	json.NewDecoder(resp.Body).Decode(&token)

	return token, nil
}

func fillTemplate(releaseBody github.ReleaseRequestBody, submitBody *SubmitBody) error {
	t1, err := template.New("t1").Parse(submitBody.Title)
	if err != nil {
		return errors.New("Couldn't parse post title: " + err.Error())
	}

	var bufTitle bytes.Buffer
	err = t1.Execute(&bufTitle, releaseBody)
	if err != nil {
		return errors.New("Couldn't parse post title: " + err.Error())
	}

	submitBody.Title = bufTitle.String()

	t2, err := template.New("t2").Parse(submitBody.Text)
	if err != nil {
		return errors.New("Couldn't parse post text: " + err.Error())
	}

	var bufText bytes.Buffer
	err = t2.Execute(&bufText, releaseBody)
	if err != nil {
		return errors.New("Couldn't parse post text: " + err.Error())
	}

	submitBody.Text = bufText.String()

	return nil

}

func Submit(token Token, releaseBody github.ReleaseRequestBody, submitBody SubmitBody) error {
	baseURL := "https://oauth.reddit.com/api/submit.json"
	submitBody.APIType = "json"

	fillTemplate(releaseBody, &submitBody)

	v := url.Values{}
	err := schema.NewEncoder().Encode(submitBody, v)
	if err != nil {
		return errors.New("Couldn't encode Reddit submission: " + err.Error())
	}

	url := baseURL + "?" + v.Encode()

	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("User-Agent", "OSS Bot 0.1")
    req.Header.Add("Content-Type", "application/x-www-form-url-encoded")
	req.Header.Add("Authorization", "Bearer " + token.AccessToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
		return errors.New("Error with Reddit submission: " + err.Error())
    }
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(bodyBytes))

    if resp.StatusCode != 200 {
		return errors.New("Unexpected status code while posting to Reddit: " + resp.Status)
    }

	return nil
}
