package translate

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	GrantType    string
	ScopeUrl     string
	ClientId     string
	ClientSecret string
	AuthUrl      string
}

type Token struct {
	AccessToken string    `json:"access_token"`
	Timestamp   time.Time `json:"-"`
	ExpiresIn   string    `json:"expires_in"`
}

func GetToken(c *Config) (token *Token, err error) {
	values := make(url.Values)
	values.Set("grant_type", c.GrantType)
	values.Set("scope", c.ScopeUrl)
	values.Set("client_id", c.ClientId)
	values.Set("client_secret", c.ClientSecret)

	resp, err := http.PostForm(c.AuthUrl, values)
	if err != nil {
		log.Fatalf("getTokens: %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll((*resp).Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New((*resp).Status)
	}
	json.Unmarshal(respBody, &token)
	token.Timestamp = time.Now()
	return
}

func (token *Token) Translate(text, from, to string) (result string, err error) {
	window, err := time.ParseDuration(token.ExpiresIn + "s")
	if err != nil {
		return "", err
	}
	if token.Timestamp.Add(window).Before(time.Now()) {
		return "", errors.New("Access token expired")
	}
	if text == "" {
		return "", errors.New("\"text\" is a required parameter")
	}
	if to == "" {
		return "", errors.New("\"to\" is a required parameter")
	}
	params := "from=" + from + "&to=" + to + "&text=" + url.QueryEscape(text)
	uri := "http://api.microsofttranslator.com/v2/Http.svc/Translate?" + params
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Content-Type", "text/plain")
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll((*resp).Body)
	err = xml.Unmarshal(bytes, &result)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", errors.New((*resp).Status)
	}
	return
}

func (token *Token) TranslateArray(texts []string, from, to string) (result []string, err error) {
	window, err := time.ParseDuration(string(token.ExpiresIn) + "s")
	if err != nil {
		return nil, err
	}
	if token.Timestamp.Add(window).Before(token.Timestamp.UTC()) {
		return nil, errors.New("Access token expired")
	}

	type Request struct {
		XMLName xml.Name `xml:"TranslateArrayRequest"`
		AppId   string
		From    string
		To      string
		Texts   []string `xml:"Texts>string"`
	}

	data, err := xml.Marshal(&Request{From: from, To: to, Texts: texts})
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(data)

	uri := "http://api.microsofttranslator.com/v2/Http.svc/TranslateArray"
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Add("Bearer", token.AccessToken)
	req.Header.Add("Content-Type", "text/xml")
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll((*resp).Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New((*resp).Status)
	}

	type TranslateArrayResponse struct {
		Error                     string
		OriginalSentenceLengths   []int
		TranslatedText            string
		TranslatedSentenceLengths []int
		State                     string
	}
	type Response struct {
		XMLName   xml.Name                 `xml:"ArrayOfTranslateArrayResponse"`
		Responses []TranslateArrayResponse `xml:"TranslateArrayResponse"`
	}

	response := Response{}
	err = xml.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	texts = []string{}
	for _, result := range response.Responses {
		texts = append(texts, result.TranslatedText)
	}

	return
}
