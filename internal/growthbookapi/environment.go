package growthbookapi

import (
	"context"
)

// CreateEnvironment creates a new environment in GrowthBook.
func (c *Client) CreateEnvironment(ctx context.Context, e *Environment) (*Environment, error) {
	out, err := fetcher[Environment](c, "POST", "/environments").One(ctx, e, "environment")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateEnvironment updates an existing environment by its ID.
func (c *Client) UpdateEnvironment(ctx context.Context, id string, e *Environment) (*Environment, error) {
	out, err := fetcher[Environment](c, "PUT", "/environments/"+id).One(ctx, e, "environment")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteEnvironment deletes an environment by its ID.
func (c *Client) DeleteEnvironment(ctx context.Context, id string) error {
	return c.delete(ctx, "/environments/"+id)
}

// FindEnvironmentByID fetches an environment by its ID by listing all and filtering.
func (c *Client) FindEnvironmentByID(ctx context.Context, id string) (*Environment, error) {
	envs, err := fetcher[[]Environment](c, "GET", "/environments").One(ctx, nil, "environments")
	if err != nil {
		return nil, err
	}
	for _, env := range envs {
		if env.ID == id {
			return &env, nil
		}
	}
	return nil, ErrNotFound
}
