//nolint:dupl

package growthbookapi

import (
	"context"
)

// CreateSDKConnection creates a new SDK connection in GrowthBook.
func (c *Client) CreateSDKConnection(ctx context.Context, s *SDKConnection) (*SDKConnection, error) {
	out, err := fetcher[SDKConnection](c, "POST", "/sdk-connections").One(ctx, s, "sdkConnection")
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
	out, err := fetcher[SDKConnection](c, "GET", "/sdk-connections/"+id).One(ctx, nil, "sdkConnection")
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
	out, err := fetcher[SDKConnection](c, "PUT", "/sdk-connections/"+id).One(ctx, s, "sdkConnection")
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
	return c.delete(ctx, "/sdk-connections/"+id)
}

// FindSDKConnectionByName searches for an SDK connection by its name and returns the first match, handling pagination.
func (c *Client) FindSDKConnectionByName(ctx context.Context, name string) (*SDKConnection, error) {
	sdks, err := fetcher[SDKConnection](c, "GET", "/sdk-connections").All(ctx, nil, "connections")
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
