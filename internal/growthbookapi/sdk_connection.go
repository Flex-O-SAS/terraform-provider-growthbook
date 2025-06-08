//nolint:dupl

package growthbookapi

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type sdkConnectionResponse struct {
	SDKConnection SDKConnection `json:"sdkConnection"`
}
type sdkConnectionListResponse struct {
	SDKConnections []SDKConnection `json:"connections"`
}

// CreateSDKConnection creates a new SDK connection in GrowthBook.
func (c *Client) CreateSDKConnection(ctx context.Context, s *SDKConnection) (*SDKConnection, error) {
	out, err := doRequestAndDecode[sdkConnectionResponse](
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
	out, err := doRequestAndDecode[sdkConnectionResponse](ctx, c, "GET", "/sdk-connections/"+id, nil, http.StatusOK)
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
	out, err := doRequestAndDecode[sdkConnectionResponse](ctx, c, "PUT", "/sdk-connections/"+id, s, http.StatusOK)
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
	return c.doDeleteRequest(ctx, "/sdk-connections/"+id, http.StatusOK, http.StatusNoContent)
}

// FindSDKConnectionByName searches for an SDK connection by its name and returns the first match.
func (c *Client) FindSDKConnectionByName(ctx context.Context, name string) (*SDKConnection, error) {
	tflog.Debug(ctx,
		"searching for sdk-connection by name",
		map[string]interface{}{
			"name": name,
		},
	)

	sdks, err := doRequestAndDecode[sdkConnectionListResponse](ctx, c, "GET", "/sdk-connections", nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for _, s := range sdks.SDKConnections {
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
