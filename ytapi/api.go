package ytapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	host          = "https://translate.yandex.net/api/v1.5/tr.json/"
	langsPath     = "getLangs"
	detectPath    = "detect"
	translagePath = "translate"
)

//var apiKey = os.Getenv("Y_TRANSLATE_API")
var apiKey string

// NewClient returns api client
func NewClient(key string) *Endpoint {
	// ignore invalid key situation
	// LOL
	return &Endpoint{key}
}

type list struct {
	Langs map[string]string `json:"langs"`
}

type detect struct {
	Code int    `json:"code"`
	Lang string `json:"lang"`
}

type translate struct {
	Detected map[string]string `json:"detected"`
	To       string            `json:"lang"` // "en-zh" or "zh-en"
	Text     []string          `json:"text"`
	Code     int               `json:"code"`
}

// Endpoint is the request endpoint
type Endpoint struct {
	apiKey string
}

// GetList get supported languages
func (e *Endpoint) GetList(ui string) (outputs map[string]string, err error) {
	params := map[string]string{
		"key": e.apiKey,
		"ui":  ui,
	}
	var result list
	f := func(r *http.Response) error {
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			return err
		}
		if len(result.Langs) == 0 {
			return fmt.Errorf("cannot find supported language for %s", ui)
		}
		outputs = result.Langs
		return nil
	}
	err = postData(host+langsPath, f, params)
	if err != nil {
		outputs = nil
	}
	return
}

// Detect the source language
func (e *Endpoint) Detect(text string) (outputs string, err error) {
	params := map[string]string{
		"key":  e.apiKey,
		"text": text,
		"hint": "en,ja,zh,de",
	}
	var result detect
	f := func(r *http.Response) error {
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			return err
		}
		if result.Code != 200 {
			return fmt.Errorf("cannot detect %s language", text)
		}
		return nil
	}
	err = postData(host+detectPath, f, params)
	if err == nil {
		outputs = result.Lang
	}
	return
}

// Trans do the translation work
func (e *Endpoint) Trans(text, to string) (outputs string, err error) {
	params := map[string]string{
		"key":     e.apiKey,
		"text":    text,
		"lang":    to, // ru-en  en
		"format":  "text",
		"options": "1",
	}
	var result translate
	f := func(r *http.Response) error {
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			return err
		}
		if r.StatusCode != 200 {
			return fmt.Errorf("something went wrong, when processing %s", text)
		}
		return nil
	}
	err = postData(host+translagePath, f, params)
	if err == nil {
		outputs = strings.Join(result.Text, "\n")
	}
	return
}

func postData(u string, fn func(*http.Response) error, params map[string]string) (err error) {
	p := url.Values{}
	for k, v := range params {
		p.Set(k, v)
	}
	resp, err := http.PostForm(u, p)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = fn(resp)
	return
}
