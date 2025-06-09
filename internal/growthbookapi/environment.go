package growthbookapi

import (
	"context"
	"net/http"
)

type environmentResponse struct {
	Environment Environment `json:"environment"`
}
type environmentListResponse struct {
	Environments []Environment `json:"environments"`
}

// CreateEnvironment creates a new environment in GrowthBook.
func (c *Client) CreateEnvironment(ctx context.Context, e *Environment) (*Environment, error) {
	out, err := fetchOne[environmentResponse](
		ctx,
		c,
		"POST",
		"/environments",
		e,
		http.StatusOK,
		http.StatusCreated,
	)
	if err != nil {
		return nil, err
	}
	return &out.Environment, nil
}

// UpdateEnvironment updates an existing environment by its ID.
func (c *Client) UpdateEnvironment(ctx context.Context, id string, e *Environment) (*Environment, error) {
	out, err := fetchOne[environmentResponse](ctx, c, "PUT", "/environments/"+id, e, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &out.Environment, nil
}

// DeleteEnvironment deletes an environment by its ID.
func (c *Client) DeleteEnvironment(ctx context.Context, id string) error {
	return c.remove(ctx, "/environments/"+id, http.StatusOK, http.StatusNoContent)
}

// FindEnvironmentByID fetches an environment by its ID by listing all and filtering.
func (c *Client) FindEnvironmentByID(ctx context.Context, id string) (*Environment, error) {
	envs, err := fetchOne[environmentListResponse](ctx, c, "GET", "/environments", nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for _, env := range envs.Environments {
		if env.ID == id {
			return &env, nil
		}
	}
	return nil, ErrNotFound
}
