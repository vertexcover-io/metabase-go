package metabase_client

import (
	"github.com/pkg/errors"
)

type DbEngineSetup struct {
	Engine         string            `json:"engine"`
	Name           string            `json:"name"`
	Details        map[string]string `json:"details"`
	IsFullSync     *bool             `json:"is_full_sync,omitempty"`
	AutoRunQueries *bool             `json:"auto_run_queries,omitempty"`
}

type UserSetup struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type PrefSetup struct {
	SiteName      string `json:"site_name"`
	SiteLocale    string `json:"site_locale,omitempty"`
	AllowTracking *bool  `json:"allow_tracking,omitempty"`
}

type InstanceSetup struct {
	User     UserSetup     `json:"user"`
	Token    string        `json:"token"`
	Database DbEngineSetup `json:"database"`
	Pref     PrefSetup     `json:"prefs"`
}

func (c *APIClient) SetupInstance(body *InstanceSetup, saveToken bool) (string, error) {
	var idResp sessionIdResponse
	if _, err := responseHandler(
		c.client.R().
			SetBody(body).
			SetResult(&idResp).
			Post(c.makeURL("setup")),
	); err != nil {
		return "", errors.Wrapf(err, "Failed to setup instance")
	}
	if saveToken {
		c = c.WithSessionToken(idResp.Id)
	}
	return idResp.Id, nil
}

func (c *APIClient) UpdateSetting(key string, value interface{}) error {
	body := map[string]interface{}{
		"value": value,
	}
	_, err := responseHandler(
		c.client.R().SetBody(body).Put(c.makeURLWithParams("setting/%s", key)),
	)
	if err != nil {
		return errors.Wrapf(err, "Unable to update setting %s=%v", key, value)
	}
	return nil
}

func (c *APIClient) EnableEmbedding(enable bool) error {
	return c.UpdateSetting("enable-embedding", enable)
}

func (c *APIClient) HealthCheck() error {
	resp, err := c.client.R().Get(c.makeURL("health"))
	if err != nil {
		return errors.Wrap(err, "Connection Error")
	}
	if resp.StatusCode() >= 500 {
		return errors.Wrap(err, "Health Check Failing")
	}
	return nil
}

func (c *APIClient) SetEmbeddingToken(token string) error {
	return c.UpdateSetting("embedding-secret-key", token)
}
