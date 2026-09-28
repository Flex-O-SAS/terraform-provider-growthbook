package growthbookapi

import "context"

func (c *Client) CreateFactMetric(ctx context.Context, fm *FactMetric) (*FactMetric, error) {
	out, err := fetcher[FactMetric](c, "POST", "/fact-metrics").One(ctx, fm, "factMetric")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetFactMetric(ctx context.Context, id string) (*FactMetric, error) {
	out, err := fetcher[FactMetric](c, "GET", "/fact-metrics/"+id).One(ctx, nil, "factMetric")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateFactMetric(ctx context.Context, id string, fm *FactMetric) (*FactMetric, error) {
	out, err := fetcher[FactMetric](c, "POST", "/fact-metrics/"+id).One(ctx, fm, "factMetric")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteFactMetric(ctx context.Context, id string) error {
	return c.delete(ctx, "/fact-metrics/"+id)
}

func (c *Client) FindFactMetricByName(ctx context.Context, name string) (*FactMetric, error) {
	items, err := fetcher[FactMetric](c, "GET", "/fact-metrics").All(ctx, nil, "factMetrics")
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
