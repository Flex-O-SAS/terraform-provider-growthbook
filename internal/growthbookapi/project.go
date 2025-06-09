//nolint:dupl

package growthbookapi

import (
	"context"
	"net/http"
)

type projectResponse struct {
	Project Project `json:"project"`
}

// CreateProject creates a new project in GrowthBook.
func (c *Client) CreateProject(ctx context.Context, p *Project) (*Project, error) {
	out, err := fetchOne[projectResponse](ctx, c, "POST", "/projects", p, http.StatusOK, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	return &out.Project, nil
}

// GetProject fetches a project by its ID.
func (c *Client) GetProject(ctx context.Context, id string) (*Project, error) {
	out, err := fetchOne[projectResponse](ctx, c, "GET", "/projects/"+id, nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &out.Project, nil
}

// UpdateProject updates an existing project by its ID.
func (c *Client) UpdateProject(ctx context.Context, id string, p *Project) (*Project, error) {
	out, err := fetchOne[projectResponse](ctx, c, "PUT", "/projects/"+id, p, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &out.Project, nil
}

// DeleteProject deletes a project by its ID.
func (c *Client) DeleteProject(ctx context.Context, id string) error {
	return c.remove(ctx, "/projects/"+id, http.StatusOK, http.StatusNoContent)
}

// FindProjectByName searches for a project by its name and returns the first match, handling pagination.
func (c *Client) FindProjectByName(ctx context.Context, name string) (*Project, error) {
	projects, err := fetchAllPages[Project](ctx, c, "GET", "/projects", nil, "projects", http.StatusOK)
	if err != nil {
		return nil, err
	}
	for _, p := range projects {
		if p.Name == name {
			return &p, nil
		}
	}
	return nil, ErrNotFound
}
