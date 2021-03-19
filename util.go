package metabase_client

import "github.com/pkg/errors"

func (c *APIClient) RandomToken() (string, error) {
	var token = struct {
		Token string `json:"token"`
	}{}
	_, err := responseHandler(
		c.client.R().SetResult(&token).Get(c.makeURL("util/random_token")),
	)
	if err != nil {
		return "", errors.Wrap(err, "Unable to generate random token")
	}

	return token.Token, err
}
