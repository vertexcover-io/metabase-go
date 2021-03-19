package metabase_client

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type visualizationSettings map[string]interface{}

type seriesSettings struct {
	MarkingEnabled bool   `mapstructure:"line.marker_enabled,omitempty"`
	MissingValue   string `mapstructure:"line.missing,omitempty"`
}

type LineVisualizationSetting struct {
	Dimension     []string       `mapstructure:"graph.dimensions,omitempty"`
	Metrics       []string       `mapstructure:"graph.metrics,omitempty"`
	XAxis         string         `mapstructure:"graph.x_axis.title_text,omitempty"`
	YAxis         string         `mapstructure:"graph.y_axis.title_text,omitempty"`
	PivotColumn   string         `mapstructure:"table.pivot_column,omitempty"`
	CellColumn    string         `mapstructure:"table.cell_column,omitempty"`
	SeriesSetting seriesSettings `mapstructure:"series_settings,omitempty"`
}

type nativeDatabaseQuery struct {
	Query        string                 `json:"query"`
	TemplateTags map[string]interface{} `json:"template_tags"`
}

type datasetQuery struct {
	DatabaseId  int                  `json:"database"`
	NativeQuery *nativeDatabaseQuery `json:"native,omitempty"`
	Type        string               `json:"type"`
}

type Card struct {
	CommonFields
	Display              string                `json:"display"`
	CollectionId         *CollectionId         `json:"collection_id"`
	VisualizationSetting visualizationSettings `json:"visualization_settings"`
	DatabaseQuery        datasetQuery          `json:"dataset_query"`
	Archived             *bool                 `json:"archived,omitempty"`
	EnableEmbedding      *bool                 `json:"enable_embedding,omitempty"`
}

type Snippet struct {
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	Id           int     `json:"id,omitempty"`
	CollectionId *int    `json:"collection_id"`
	Content      string  `json:"content"`
}

type CardOption func(card *Card) (*Card, error)

func NativeQuery(databaseId int, query string, templateTags map[string]interface{}) CardOption {
	return func(card *Card) (*Card, error) {
		card.DatabaseQuery = datasetQuery{
			DatabaseId: databaseId,
			NativeQuery: &nativeDatabaseQuery{
				Query:        query,
				TemplateTags: templateTags,
			},
			Type: "native",
		}
		return card, nil
	}
}

func SetDisplay(display string, settings map[string]interface{}) CardOption {
	return func(card *Card) (*Card, error) {
		if settings == nil {
			settings = make(map[string]interface{})
		}
		card.Display = display
		card.VisualizationSetting = settings
		return card, nil
	}
}

func SetDisplayAsScalar() CardOption {
	return SetDisplay("scalar", map[string]interface{}{})
}

func SetDisplayAsLine(settings *LineVisualizationSetting) CardOption {
	return func(card *Card) (*Card, error) {
		card.Display = "line"
		var visSettings visualizationSettings
		if err := mapstructure.Decode(settings, visSettings); err != nil {
			return card, errors.Wrap(err, "Unable to parse visualization settings")
		}
		card.VisualizationSetting = visSettings
		return card, nil
	}
}

func SetVisualizationSettings(settings map[string]interface{}) CardOption {
	return func(card *Card) (*Card, error) {
		card.VisualizationSetting = settings
		return card, nil
	}
}

func AllowCardEmbedding() CardOption {
	return func(card *Card) (*Card, error) {
		value := true
		card.EnableEmbedding = &value
		return card, nil
	}
}

func (c *APIClient) CreateCard(
	name string,
	description *string,
	collectionId *CollectionId,
	cardOptions ...CardOption) (*Card, error) {
	card := &Card{
		CommonFields: CommonFields{
			Name:        name,
			Description: description,
		},
		CollectionId: collectionId,
	}
	for _, option := range cardOptions {
		var err error
		card, err = option(card)
		if err != nil {
			return card, errors.Wrap(err, "Error while constructing question card")
		}
	}
	if _, err := responseHandler(
		c.client.R().SetResult(card).
			SetBody(card).
			Post(c.makeURL("card")),
	); err != nil {
		return nil, errors.Wrapf(err, "Unable to create card: %s", name)
	}
	return card, nil
}

func (c *APIClient) CreateSnippet(snippet *Snippet) (*Snippet, error) {
	if _, err := c.client.R().SetResult(snippet).
		SetBody(snippet).Post("card"); err != nil {
		return nil, errors.Wrap(err, "Unable to create card")
	}
	return snippet, nil
}

func (c *APIClient) GetAllCards(
	filter *string,
) ([]Card, error) {
	var cards []Card = make([]Card, 0)
	qs := make(map[string]string)
	if filter != nil {
		qs["f"] = *filter
	}
	_, err := responseHandler(
		c.client.R().
			SetResult(&cards).SetQueryParams(qs).
			Get(c.makeURL("card")),
	)
	return cards, errors.Wrap(err, "Unable to fetch list of cards")
}
