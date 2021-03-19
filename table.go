package metabase_client

import "github.com/pkg/errors"

type Table struct {
	CommonFields
	DbID        int    `json:"db_id"`
	EntityType  string `json:"entity_type"`
	Schema      string `json:"schema"`
	DisplayName string `json:"display_name"`
	Active      *bool  `json:"active"`
}
type TableWithMetadata struct {
	Table
	Fields []TableField `json:"fields"`
}

type TableField struct {
	CommonFields
	DatabaseType           string                 `json:"database_type"`
	TableID                int                    `json:"table_id"`
	SpecialType            *string                `json:"special_type"`
	FingerprintVersion     int                    `json:"fingerprint_version"`
	HasFieldValues         string                 `json:"has_field_values"`
	FkTargetFieldID        *int                   `json:"fk_target_field_id"`
	Dimensions             []string               `json:"dimensions"`
	DimensionOptions       []string               `json:"dimension_options"`
	CustomPosition         *int                   `json:"custom_position"`
	Active                 bool                   `json:"active"`
	ParentID               *int                   `json:"parent_id"`
	LastAnalyzed           *string                `json:"last_analyzed"`
	Position               *int                   `json:"position"`
	VisibilityType         *string                `json:"visibility_type"`
	DefaultDimensionOption *string                `json:"default_dimension_option"`
	PreviewDisplay         bool                   `json:"preview_display"`
	DisplayName            string                 `json:"display_name"`
	DatabasePosition       int                    `json:"database_position"`
	Fingerprint            map[string]interface{} `json:"fingerprint"`
	BaseType               *string                `json:"base_type"`
}

func (c *APIClient) GetTableMetadata(tableId int, includeHiddenFields bool, includeSensitiveFields bool) (TableWithMetadata, error) {
	metadata := TableWithMetadata{}
	qs := make(map[string]string)
	if includeHiddenFields {
		qs["include_hidden_fields"] = "true"
	}
	if includeSensitiveFields {
		qs["include_sensitive_fields"] = "true"
	}
	_, err := responseHandler(
		c.client.R().
			SetResult(&metadata).SetQueryParams(qs).
			Get(c.makeURLWithParams("table/%d/query_metadata", tableId)),
	)
	return metadata, errors.Wrapf(err, "Unable to fetch table metadata for: %d", tableId)
}
