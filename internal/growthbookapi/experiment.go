package growthbookapi

import "context"

func (c *Client) CreateExperiment(ctx context.Context, e *Experiment) (*Experiment, error) {
	out, err := fetcher[Experiment](c, "POST", "/experiments").One(ctx, e, "experiment")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetExperiment(ctx context.Context, id string) (*Experiment, error) {
	out, err := fetcher[Experiment](c, "GET", "/experiments/"+id).One(ctx, nil, "experiment")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateExperiment(ctx context.Context, id string, e *Experiment) (*Experiment, error) {
	out, err := fetcher[Experiment](c, "POST", "/experiments/"+id).One(ctx, e, "experiment")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteExperiment(ctx context.Context, id string) error {
	return c.delete(ctx, "/experiments/"+id)
}

func (c *Client) FindExperimentByName(ctx context.Context, name string) (*Experiment, error) {
	items, err := fetcher[Experiment](c, "GET", "/experiments").All(ctx, nil, "experiments")
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
