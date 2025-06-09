//nolint:dupl

package growthbookapi

import (
	"context"
	"net/http"
)

type sdkConnectionResponse struct {
	SDKConnection SDKConnection `json:"sdkConnection"`
}

// CreateSDKConnection creates a new SDK connection in GrowthBook.
func (c *Client) CreateSDKConnection(ctx context.Context, s *SDKConnection) (*SDKConnection, error) {
	out, err := fetchOne[sdkConnectionResponse](
		ctx,
		c,
		"POST",
		"/sdk-connections",
		s,
		http.StatusOK,
		http.StatusCreated,
	)
	if err != nil {
		return nil, err
	}
	if len(out.SDKConnection.Languages) != 0 {
		out.SDKConnection.Language = out.SDKConnection.Languages[0]
	}
	return &out.SDKConnection, nil
}

// GetSDKConnection fetches an SDK connection by its ID.
func (c *Client) GetSDKConnection(ctx context.Context, id string) (*SDKConnection, error) {
	out, err := fetchOne[sdkConnectionResponse](ctx, c, "GET", "/sdk-connections/"+id, nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	if len(out.SDKConnection.Languages) != 0 {
		out.SDKConnection.Language = out.SDKConnection.Languages[0]
	}
	return &out.SDKConnection, nil
}

// UpdateSDKConnection updates an existing SDK connection by its ID.
func (c *Client) UpdateSDKConnection(ctx context.Context, id string, s *SDKConnection) (*SDKConnection, error) {
	out, err := fetchOne[sdkConnectionResponse](ctx, c, "PUT", "/sdk-connections/"+id, s, http.StatusOK)
	if err != nil {
		return nil, err
	}
	if len(out.SDKConnection.Languages) != 0 {
		out.SDKConnection.Language = out.SDKConnection.Languages[0]
	}
	return &out.SDKConnection, nil
}

// DeleteSDKConnection deletes an SDK connection by its ID.
func (c *Client) DeleteSDKConnection(ctx context.Context, id string) error {
	return c.remove(ctx, "/sdk-connections/"+id, http.StatusOK, http.StatusNoContent)
}

// FindSDKConnectionByName searches for an SDK connection by its name and returns the first match, handling pagination.
func (c *Client) FindSDKConnectionByName(ctx context.Context, name string) (*SDKConnection, error) {
	sdkConnections, err := fetchAllPages[SDKConnection](
		ctx,
		c,
		"GET",
		"/sdk-connections",
		nil,
		"connections",
		http.StatusOK,
	)
	if err != nil {
		return nil, err
	}
	for _, s := range sdkConnections {
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
