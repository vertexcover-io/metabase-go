package metabase_client

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	authHeader = "X-Metabase-Session"
)

type APIClient struct {
	BasePath string `json:"basePath,omitempty"`
	client   *resty.Client
}

func NewAPIClient(host string) *APIClient {
	return &APIClient{
		BasePath: fmt.Sprintf("%s/api", host),
		client: resty.New().
			SetHeader("Accept", "application/json").
			SetHeader("Content-Type", "application/json"),
	}
}

func (c *APIClient) WithDefaultHeader(key string, val string) *APIClient {
	c.client = c.client.SetHeader(key, val)
	return c
}

func (c *APIClient) WithSessionToken(token string) *APIClient {
	return c.WithDefaultHeader(authHeader, token)
}

func (c *APIClient) makeURL(path string) string {
	return fmt.Sprintf("%s/%s", c.BasePath, path)
}

func (c *APIClient) makeURLWithParams(path string, params ...interface{}) string {
	return c.makeURL(fmt.Sprintf(path, params...))
}
