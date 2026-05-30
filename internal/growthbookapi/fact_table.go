package growthbookapi

import "context"

func (c *Client) CreateFactTable(ctx context.Context, ft *FactTable) (*FactTable, error) {
	out, err := fetcher[FactTable](c, "POST", "/fact-tables").One(ctx, ft, "factTable")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetFactTable(ctx context.Context, id string) (*FactTable, error) {
	out, err := fetcher[FactTable](c, "GET", "/fact-tables/"+id).One(ctx, nil, "factTable")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateFactTable(ctx context.Context, id string, ft *FactTable) (*FactTable, error) {
	out, err := fetcher[FactTable](c, "POST", "/fact-tables/"+id).One(ctx, ft, "factTable")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteFactTable(ctx context.Context, id string) error {
	return c.delete(ctx, "/fact-tables/"+id)
}

func (c *Client) FindFactTableByName(ctx context.Context, name string) (*FactTable, error) {
	items, err := fetcher[FactTable](c, "GET", "/fact-tables").All(ctx, nil, "factTables")
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
