package growthbookapi

import (
	"context"
)

// CreateTeam creates a new team in GrowthBook.
func (c *Client) CreateTeam(ctx context.Context, t *Team) (*Team, error) {
	out, err := fetcher[Team](c, "POST", "/teams").One(ctx, t, "team")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTeam fetches a team by its ID.
func (c *Client) GetTeam(ctx context.Context, id string) (*Team, error) {
	out, err := fetcher[Team](c, "GET", "/teams/"+id).One(ctx, nil, "team")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateTeam updates an existing team by its ID.
func (c *Client) UpdateTeam(ctx context.Context, id string, t *Team) (*Team, error) {
	out, err := fetcher[Team](c, "PUT", "/teams/"+id).One(ctx, t, "team")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteTeam deletes a team by its ID.
func (c *Client) DeleteTeam(ctx context.Context, id string) error {
	return c.delete(ctx, "/teams/"+id)
}

// FindTeamByName searches for a team by its name and returns the first match.
func (c *Client) FindTeamByName(ctx context.Context, name string) (*Team, error) {
	teams, err := fetcher[Team](c, "GET", "/teams").All(ctx, nil, "teams")
	if err != nil {
		return nil, err
	}
	for _, t := range teams {
		if t.Name == name {
			return &t, nil
		}
	}
	return nil, ErrNotFound
}
