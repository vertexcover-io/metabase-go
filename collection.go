package metabase_client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const roolCollectionID = -1

type CollectionId struct {
	int
}

func (c CollectionId) MarshalJSON() ([]byte, error) {
	if c.int == roolCollectionID {
		return []byte("root"), nil
	}
	return json.Marshal(c.int)
}

func (c *CollectionId) UnmarshalJSON(b []byte) error {
	idAsStr := strings.Trim(string(b), "\"")
	if idAsStr == "root" {
		c.int = roolCollectionID
		return nil
	} else {
		c.int = 0
		return json.Unmarshal(b, &(c.int))
	}
}

func (c *CollectionId) Int() int {
	return c.int
}

type Collection struct {
	Id          CollectionId  `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	Color       *string       `json:"color,omitempty"`
	Description *string       `json:"description,omitempty"`
	Namespace   *string       `json:"namespace,omitempty"`
	Archived    *bool         `json:"archived,omitempty"`
	ParentId    *CollectionId `json:"parent_id,omitempty"`
}

func (c *Collection) GetUrl() string {
	return fmt.Sprintf("/collection/%d", c.Id.int)
}

type CollectionItem struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Favorite    bool    `json:"favorite"`
	Model       string  `json:"model"`
}

func (c *APIClient) GetAllCollections(isArchived *bool, namespace *string) ([]Collection, error) {
	var collections []Collection = make([]Collection, 0)
	qs := make(map[string]string)
	if isArchived != nil {
		if *isArchived {
			qs["archived"] = "true"
		} else {
			qs["archived"] = "false"
		}
	}
	if namespace != nil {
		qs["namepspace"] = *namespace
	}
	_, err := c.client.R().
		SetResult(&collections).SetQueryParams(qs).
		Get(c.makeURL("collection"))
	return collections[1:], errors.Wrap(err, "Unable to fetch collections")
}

func setCollectionFields(
	collection *Collection, description *string, color *string, parent_id *CollectionId, namespace *string,
	archived *bool,
) *Collection {
	if color != nil {
		collection.Color = color
	}

	if description != nil {
		collection.Description = description
	}

	if parent_id != nil {
		collection.ParentId = parent_id
	}

	if namespace != nil {
		collection.Namespace = namespace
	}

	if archived != nil {
		collection.Archived = archived
	}

	return collection
}

func (c *APIClient) CreateCollection(
	name string, description *string, color *string, parent_id *CollectionId, namespace *string,
) (*Collection, error) {
	collection := &Collection{
		Name: name,
	}

	if color == nil {
		c := "#007bff"
		color = &c
	}

	collection = setCollectionFields(collection, description, color, parent_id, namespace, nil)

	if _, err := responseHandler(
		c.client.R().
			SetResult(collection).
			SetBody(collection).Post(c.makeURL("collection"))); err != nil {
		return nil, errors.Wrapf(err, "Unable to create collection: %s", name)
	}
	return collection, nil
}

func (c *APIClient) UpdateCollection(
	collectionID CollectionId,
	description *string,
	color *string,
	parentId *CollectionId,
	namespace *string,
	archived *bool,
) (*Collection, error) {
	collection := &Collection{
		Id: collectionID,
	}
	collection = setCollectionFields(collection, description, color, parentId, namespace, archived)

	if _, err := responseHandler(
		c.client.R().
			SetResult(collection).
			SetBody(collection).Put(c.makeURLWithParams("collection/%d", collectionID.int))); err != nil {
		return nil, errors.Wrapf(err, "Unable to update collection: %d", collectionID)
	}

	return collection, nil
}

func (c *APIClient) GetCollectionItems(id int, model *string, isArchived *bool) ([]CollectionItem, error) {
	var collectionItems []CollectionItem = make([]CollectionItem, 0)
	qs := make(map[string]string)
	if isArchived != nil {
		if *isArchived {
			qs["archived"] = "true"
		} else {
			qs["archived"] = "false"
		}
	}

	if model != nil {
		qs["model"] = *model
	}

	_, err := c.client.R().
		SetResult(&collectionItems).SetQueryParams(qs).
		Get(c.makeURLWithParams("collection/%d/items", id))
	return collectionItems, errors.Wrap(err, "Unable to fetch collection items")
}
