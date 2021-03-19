package metabase_client

import (
	"math/rand"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func responseHandler(resp *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		return resp, err
	}
	if resp.IsSuccess() {
		return resp, nil
	} else {
		return resp, errors.Errorf("Request Failed. Status Code: %d, Response: %+v", resp.StatusCode(), string(resp.Body()))
	}
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
