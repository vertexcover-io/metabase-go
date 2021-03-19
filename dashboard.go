package metabase_client

import (
	"strings"

	"github.com/pkg/errors"
)

type ParamMapping struct {
	CardId      int           `json:"card_id"`
	ParameterId string        `json:"parameter_id"`
	Target      []interface{} `json:"target"`
}

type OrderedCard struct {
	CommonFields
	CardId                int                    `json:"card_id"`
	DashboardId           int                    `json:"dashboard_id"`
	SizeX                 int                    `json:"sizeX"`
	SizeY                 int                    `json:"sizeY"`
	Row                   int                    `json:"row"`
	Column                int                    `json:"col"`
	Series                []interface{}          `json:"series"`
	ParameterMappings     []ParamMapping         `json:"parameter_mappings"`
	VisualizationSettings map[string]interface{} `json:"visualization_settings"`
}

type FilterParams struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Slug    string  `json:"slug"`
	Type    string  `json:"type"`
	Default *string `json:"default,omitempty"`
}

type Dashboard struct {
	CommonFields
	CollectionId        *CollectionId          `json:"collection_id,omitempty"`
	CollectionPoisition *int                   `json:"collection_position,omitempty"`
	OrderedCards        []OrderedCard          `json:"ordered_cards"`
	Archived            *bool                  `json:"archived,omitempty"`
	CanWrite            *bool                  `json:"can_write,omitempty"`
	EnableEmbedding     *bool                  `json:"enable_embedding,omitempty"`
	EmbeddingParams     map[string]interface{} `json:"embedding_params,omitempty"`
	CreatorID           *int                   `json:"creator_id,omitempty"`
	ParamFields         map[string]interface{} `json:"param_fields,omitempty"`
	ParamValues         interface{}            `json:"param_values,omitempty"`
	Parameters          []FilterParams         `json:"parameters,omitempty"`
}

type DashboardOption func(dashboard *Dashboard) (*Dashboard, error)

func WithParameter(
	param FilterParams,
) DashboardOption {
	return func(dashboard *Dashboard) (*Dashboard, error) {
		dashboard.Parameters = append(dashboard.Parameters, param)
		return dashboard, nil
	}
}

func ToggleDashboardEmbedding(enable bool, params map[string]interface{}) DashboardOption {
	return func(dashboard *Dashboard) (*Dashboard, error) {
		value := true
		dashboard.EnableEmbedding = &value
		dashboard.EmbeddingParams = params
		return dashboard, nil
	}
}

func MakeFilterParam(name string, type_ string, default_ *string) FilterParams {
	return FilterParams{
		Id:      randomString(8),
		Name:    name,
		Slug:    strings.ToLower(name),
		Type:    type_,
		Default: default_,
	}
}

func (c *APIClient) CreateDashboard(
	name string,
	description *string,
	collectionId *CollectionId,
	collectionPosition *int,
	parameters []FilterParams,
) (*Dashboard, error) {

	dashboard := &Dashboard{
		CommonFields: CommonFields{
			Name:        name,
			Description: description,
		},
		CollectionId:        collectionId,
		CollectionPoisition: collectionPosition,
		Parameters:          parameters,
	}
	if _, err := responseHandler(
		c.client.R().SetResult(dashboard).
			SetBody(dashboard).
			Post(c.makeURL("dashboard")),
	); err != nil {
		return nil, errors.Wrapf(err, "Unable to create dashboard: %s", name)
	}
	return dashboard, nil
}

func (c *APIClient) AddCardToDashboard(
	dashboardId int,
	cardId int,
) (OrderedCard, error) {
	body := map[string]interface{}{
		"cardId": cardId,
	}
	card := OrderedCard{}

	if _, err := responseHandler(
		c.client.R().SetResult(&card).
			SetBody(body).
			Post(c.makeURLWithParams("dashboard/%d/cards", dashboardId)),
	); err != nil {
		return card, errors.Wrapf(err, "Unable to update dashboard: %d", dashboardId)
	}
	return card, nil
}

func (c *APIClient) UpdateDashboard(
	dashboardId int,
	dashboardOptions ...DashboardOption,
) (*Dashboard, error) {
	dashboard := &Dashboard{
		CommonFields: CommonFields{
			Id: dashboardId,
		},
	}
	for _, option := range dashboardOptions {
		var err error
		dashboard, err = option(dashboard)
		if err != nil {
			return dashboard, errors.Wrap(err, "Error while constructing dashboard")
		}
	}
	if _, err := responseHandler(
		c.client.R().SetResult(dashboard).
			SetBody(dashboard).
			Put(c.makeURLWithParams("dashboard/%d", dashboardId)),
	); err != nil {
		return nil, errors.Wrapf(err, "Unable to update dashboard: %d", dashboardId)
	}
	return dashboard, nil
}

func (c *APIClient) UpdateDashboardCards(
	dashboardId int,
	orderedCards ...OrderedCard,
) error {
	body := map[string]interface{}{
		"id":    dashboardId,
		"cards": orderedCards,
	}
	if _, err := responseHandler(
		c.client.R().
			SetBody(body).
			Put(c.makeURLWithParams("dashboard/%d/cards", dashboardId)),
	); err != nil {
		return errors.Wrapf(err, "Unable to update dashboard cards: %d", dashboardId)
	}
	return nil
}

func (c *APIClient) GetAllDashboards(
	filter *string,
) ([]Dashboard, error) {
	var dashboards []Dashboard = make([]Dashboard, 0)
	qs := make(map[string]string)
	if filter != nil {
		qs["f"] = *filter
	}
	_, err := responseHandler(
		c.client.R().
			SetResult(&dashboards).SetQueryParams(qs).
			Get(c.makeURL("dashboard")),
	)
	return dashboards, errors.Wrap(err, "Unable to fetch list of dashboards")
}
