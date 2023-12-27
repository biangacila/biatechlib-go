package dropbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	AUTHORIZATION_URL  = "https://www1.dropbox.com/1/oauth2/authorize?client_id=%v&response_type=code&redirect_uri=%v"
	TOKEN_EXCHANGE_URL = "https://api.dropbox.com/1/oauth2/token"
)

type HttpInterface interface {
	PostForm(string, url.Values) (*http.Response, error)
}

type OAuth2Handler struct {
	Key,
	Secret,
	RedirectURL string
	authURL string
}

func newOAuth2Handler() *OAuth2Handler {
	return &OAuth2Handler{
		authURL: TOKEN_EXCHANGE_URL,
	}
}

func (h *OAuth2Handler) setAppKeys(appKey string, appSecret string, redirectURL string) {
	h.Key = appKey
	h.Secret = appSecret
	h.RedirectURL = redirectURL
}

func (h *OAuth2Handler) authCodeURL() string {
	return fmt.Sprintf(AUTHORIZATION_URL, h.Key, h.RedirectURL)
}

func (h *OAuth2Handler) tokenExchange(code string) (*Token, error) {
	data := url.Values{}

	data.Add("code", code)
	data.Add("grant_type", "authorization_code")
	data.Add("client_id", h.Key)
	data.Add("client_secret", h.Secret)
	data.Add("redirect_uri", h.RedirectURL)

	resp, err := http.PostForm(h.authURL, data)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	token := &Token{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&token)

	return token, nil
}
