package growthbookapi

import (
	"context"
	"errors"
)

type AttributeUpdateBody struct {
	DataType    string   `json:"datatype"`
	Format      string   `json:"format"`
	EnumValues  string   `json:"enum"`
	Projects    []string `json:"projects"`
	Archived    bool     `json:"archived"`
	Description string   `json:"description"`
}

func (c *Client) CreateAttribute(ctx context.Context, a *Attribute) (*Attribute, error) {
	out, err := fetcher[Attribute](c, "POST", "/attributes").One(ctx, a, "attribute")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetAttribute(ctx context.Context, property string) (*Attribute, error) {
	out, err := fetcher[[]Attribute](c, "GET", "/attributes").One(ctx, nil, "attributes")
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(out); i++ {
		if out[i].Property == property {
			return &out[i], nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) UpdateAttribute(ctx context.Context, property string, a *Attribute) (*Attribute, error) {
	body := &AttributeUpdateBody{
		DataType:    a.DataType,
		Format:      a.Format,
		EnumValues:  a.EnumValues,
		Projects:    a.Projects,
		Archived:    a.Archived,
		Description: a.Description,
	}
	out, err := fetcher[Attribute](c, "PUT", "/attributes/"+property).One(ctx, body, "attribute")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteAttribute(ctx context.Context, property string) error {
	err := c.delete(ctx, "/attributes/"+property)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	return err
}
