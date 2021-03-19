package metabase_client

import (
	"fmt"
	"strings"
	"time"
)

const defaultLayout = "2006-01-02T15:04:05.99999"

type customTime struct {
	time.Time
}

func (ct *customTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return nil
	}
	var err error
	ct.Time, err = time.Parse(defaultLayout, s)
	return err
}

func (ct *customTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == (time.Time{}).UnixNano() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(defaultLayout))), nil
}

type CommonFields struct {
	Id          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   customTime `json:"created_at,omitmepty"`
	UpdatedAt   customTime `json:"updated_at,omitmepty"`
}
