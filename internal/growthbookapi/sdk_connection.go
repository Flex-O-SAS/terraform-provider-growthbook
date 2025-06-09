//nolint:dupl

package growthbookapi

import (
	"context"
)

// CreateSDKConnection creates a new SDK connection in GrowthBook.
func (c *Client) CreateSDKConnection(ctx context.Context, s *SDKConnection) (*SDKConnection, error) {
	out, err := fetchSingle[SDKConnection](ctx, c, "POST", "/sdk-connections", s, "sdkConnection")
	if err != nil {
		return nil, err
	}
	if len(out.Languages) != 0 {
		out.Language = out.Languages[0]
	}
	return &out, nil
}

// GetSDKConnection fetches an SDK connection by its ID.
func (c *Client) GetSDKConnection(ctx context.Context, id string) (*SDKConnection, error) {
	out, err := fetchSingle[SDKConnection](ctx, c, "GET", "/sdk-connections/"+id, nil, "sdkConnection")
	if err != nil {
		return nil, err
	}
	if len(out.Languages) != 0 {
		out.Language = out.Languages[0]
	}
	return &out, nil
}

// UpdateSDKConnection updates an existing SDK connection by its ID.
func (c *Client) UpdateSDKConnection(ctx context.Context, id string, s *SDKConnection) (*SDKConnection, error) {
	out, err := fetchSingle[SDKConnection](ctx, c, "PUT", "/sdk-connections/"+id, s, "sdkConnection")
	if err != nil {
		return nil, err
	}
	if len(out.Languages) != 0 {
		out.Language = out.Languages[0]
	}
	return &out, nil
}

// DeleteSDKConnection deletes an SDK connection by its ID.
func (c *Client) DeleteSDKConnection(ctx context.Context, id string) error {
	return c.remove(ctx, "/sdk-connections/"+id)
}

// FindSDKConnectionByName searches for an SDK connection by its name and returns the first match, handling pagination.
func (c *Client) FindSDKConnectionByName(ctx context.Context, name string) (*SDKConnection, error) {
	sdks, err := fetchAll[SDKConnection](ctx, c, "GET", "/sdk-connections", nil, "connections")
	if err != nil {
		return nil, err
	}
	for _, s := range sdks {
		if s.Name == name {
			if len(s.Languages) != 0 {
				s.Language = s.Languages[0]
			}
			return &s, nil
		}
	}
	return nil, ErrNotFound
}

// FindSDKConnectionByID fetches an SDK connection by its ID.
func (c *Client) FindSDKConnectionByID(ctx context.Context, id string) (*SDKConnection, error) {
	return c.GetSDKConnection(ctx, id)
}
