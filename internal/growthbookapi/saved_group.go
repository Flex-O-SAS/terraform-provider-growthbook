package growthbookapi

import "context"

func (c *Client) CreateSavedGroup(ctx context.Context, sg *SavedGroup) (*SavedGroup, error) {
	out, err := fetcher[SavedGroup](c, "POST", "/saved-groups").One(ctx, sg, "savedGroup")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetSavedGroup(ctx context.Context, id string) (*SavedGroup, error) {
	out, err := fetcher[SavedGroup](c, "GET", "/saved-groups/"+id).One(ctx, nil, "savedGroup")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateSavedGroup(ctx context.Context, id string, sg *SavedGroup) (*SavedGroup, error) {
	out, err := fetcher[SavedGroup](c, "POST", "/saved-groups/"+id).One(ctx, sg, "savedGroup")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteSavedGroup(ctx context.Context, id string) error {
	return c.delete(ctx, "/saved-groups/"+id)
}

func (c *Client) FindSavedGroupByName(ctx context.Context, name string) (*SavedGroup, error) {
	items, err := fetcher[SavedGroup](c, "GET", "/saved-groups").All(ctx, nil, "savedGroups")
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
