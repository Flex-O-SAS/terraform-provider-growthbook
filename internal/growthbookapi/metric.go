package growthbookapi

import "context"

func (c *Client) CreateMetric(ctx context.Context, m *Metric) (*Metric, error) {
	out, err := fetcher[Metric](c, "POST", "/metrics").One(ctx, m, "metric")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetMetric(ctx context.Context, id string) (*Metric, error) {
	out, err := fetcher[Metric](c, "GET", "/metrics/"+id).One(ctx, nil, "metric")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMetric sends a PUT and then re-fetches via GET (PUT returns {updatedId} only).
func (c *Client) UpdateMetric(ctx context.Context, id string, m *Metric) (*Metric, error) {
	_, err := fetcher[string](c, "PUT", "/metrics/"+id).One(ctx, m, "updatedId")
	if err != nil {
		return nil, err
	}
	return c.GetMetric(ctx, id)
}

func (c *Client) DeleteMetric(ctx context.Context, id string) error {
	return c.delete(ctx, "/metrics/"+id)
}

func (c *Client) FindMetricByName(ctx context.Context, name string) (*Metric, error) {
	items, err := fetcher[Metric](c, "GET", "/metrics").All(ctx, nil, "metrics")
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
