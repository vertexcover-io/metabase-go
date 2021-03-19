package metabase_client

import (
	"fmt"

	"github.com/pkg/errors"
)

type Database struct {
	CommonFields
	Timezone       string                 `json:"timezone"`
	AutoRunQueries *bool                  `json:"auto_run_queries"`
	IsFullSync     *bool                  `json:"is_full_sync"`
	Details        map[string]interface{} `json:"details"`
	Engine         string                 `json:"engine"`
	Tables         []Table                `json:"tables"`
	Features       []string               `json:"features"`
}

func (c *APIClient) GetAllDatabases(includeTables bool, includeSavedQuestions bool) ([]Database, error) {
	var databases []Database = make([]Database, 0)
	qs := make(map[string]string)
	if includeTables {
		qs["include"] = "tables"
	}
	if includeSavedQuestions {
		qs["saved"] = "true"
	}
	_, err := responseHandler(
		c.client.R().
			SetResult(&databases).SetQueryParams(qs).
			Get(c.makeURL("database")),
	)
	if err != nil {
		return databases, errors.Wrap(err, "Unable to fetch list of databases")
	}
	return databases, nil
}

func (c *APIClient) GetDatabase(id int, includeTables bool, includeTableFields bool) (*Database, error) {
	database := Database{}
	qs := make(map[string]string)

	if includeTableFields && includeTables {
		return &database, fmt.Errorf("You cannot include tables and fields together")
	} else if includeTables {
		qs["include"] = "tables"
	} else if includeTableFields {
		qs["include"] = "tables.fields"
	}
	_, err := responseHandler(
		c.client.R().
			SetResult(&database).SetQueryParams(qs).
			Get(c.makeURLWithParams("database/%d", id)),
	)
	if err != nil {
		return &database, errors.Wrap(err, "Unable to fetch database")
	}
	return &database, nil
}

func (c *APIClient) SyncSchema(databaseId int) error {
	_, err := responseHandler(
		c.client.R().
			Post(c.makeURLWithParams("database/%d/sync_schema", databaseId)),
	)
	if err != nil {
		return errors.Wrap(err, "Error while synching schema")
	}
	return nil
}
