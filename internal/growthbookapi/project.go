//nolint:dupl

package growthbookapi

import (
	"context"
)

// CreateProject creates a new project in GrowthBook.
func (c *Client) CreateProject(ctx context.Context, p *Project) (*Project, error) {
	out, err := fetcher[Project](c, "POST", "/projects").One(ctx, p, "project")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetProject fetches a project by its ID.
func (c *Client) GetProject(ctx context.Context, id string) (*Project, error) {
	out, err := fetcher[Project](c, "GET", "/projects/"+id).One(ctx, nil, "project")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateProject updates an existing project by its ID.
func (c *Client) UpdateProject(ctx context.Context, id string, p *Project) (*Project, error) {
	out, err := fetcher[Project](c, "PUT", "/projects/"+id).One(ctx, p, "project")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteProject deletes a project by its ID.
func (c *Client) DeleteProject(ctx context.Context, id string) error {
	return c.delete(ctx, "/projects/"+id)
}

// FindProjectByName searches for a project by its name and returns the first match, handling pagination.
func (c *Client) FindProjectByName(ctx context.Context, name string) (*Project, error) {
	projects, err := fetcher[Project](c, "GET", "/projects").All(ctx, nil, "projects")
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
