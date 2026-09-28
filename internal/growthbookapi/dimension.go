package growthbookapi

import "context"

func (c *Client) CreateDimension(ctx context.Context, d *Dimension) (*Dimension, error) {
	out, err := fetcher[Dimension](c, "POST", "/dimensions").One(ctx, d, "dimension")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetDimension(ctx context.Context, id string) (*Dimension, error) {
	out, err := fetcher[Dimension](c, "GET", "/dimensions/"+id).One(ctx, nil, "dimension")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateDimension(ctx context.Context, id string, d *Dimension) (*Dimension, error) {
	out, err := fetcher[Dimension](c, "POST", "/dimensions/"+id).One(ctx, d, "dimension")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteDimension(ctx context.Context, id string) error {
	return c.delete(ctx, "/dimensions/"+id)
}

func (c *Client) FindDimensionByName(ctx context.Context, name string) (*Dimension, error) {
	items, err := fetcher[Dimension](c, "GET", "/dimensions").All(ctx, nil, "dimensions")
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.Name == name {
			return &item, nil
		}
	}
	return nil, ErrNotFound
}
