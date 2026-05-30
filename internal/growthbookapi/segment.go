package growthbookapi

import "context"

func (c *Client) CreateSegment(ctx context.Context, s *Segment) (*Segment, error) {
	out, err := fetcher[Segment](c, "POST", "/segments").One(ctx, s, "segment")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetSegment(ctx context.Context, id string) (*Segment, error) {
	out, err := fetcher[Segment](c, "GET", "/segments/"+id).One(ctx, nil, "segment")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateSegment(ctx context.Context, id string, s *Segment) (*Segment, error) {
	out, err := fetcher[Segment](c, "POST", "/segments/"+id).One(ctx, s, "segment")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteSegment(ctx context.Context, id string) error {
	return c.delete(ctx, "/segments/"+id)
}

func (c *Client) FindSegmentByName(ctx context.Context, name string) (*Segment, error) {
	items, err := fetcher[Segment](c, "GET", "/segments").All(ctx, nil, "segments")
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
