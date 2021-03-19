package metabase_client

import (
	"github.com/pkg/errors"
)

type SessionProperties struct {
	GACode              string `json:"ga-code"`
	SiteLocale          string `json:"en"`
	ApplicationName     string `json:"application-name"`
	SiteUrl             string `json:"site-url"`
	EnablePasswordLogin bool   `json:"enable-password-login"`
	SetupToken          string `json:"setup-token"`
	EmbeddingSecretKey  string `json:"embedding-secret-key"`
	EnableEmbedding     bool   `json:"enable-embedding"`
	ShowHomePageXRays   bool   `json:"show-homepage-xrays"`
}

type sessionIdResponse struct {
	Id string `json:"id"`
}

func (c *APIClient) GetSessionProperties() (*SessionProperties, error) {
	var props SessionProperties
	_, err := c.client.R().
		SetResult(&props).
		Get(c.makeURL("session/properties"))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get session properties")
	}
	return &props, nil
}

func (c *APIClient) Login(username string, password string, saveToken bool) (string, error) {
	var idResp sessionIdResponse
	var body = struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{username, password}

	_, err := c.client.R().
		SetBody(body).
		SetResult(&idResp).
		Post(c.makeURL("session"))
	if err != nil {
		return "", errors.Wrapf(err, "Failed to Login")
	}

	if saveToken {
		c = c.WithSessionToken(idResp.Id)
	}
	return idResp.Id, nil
}
