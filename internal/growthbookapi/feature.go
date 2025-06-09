package growthbookapi

import (
	"context"
	"net/http"
)

type featureResponse struct {
	Feature Feature `json:"feature"`
}

// CreateFeature creates a new feature in GrowthBook.
func (c *Client) CreateFeature(ctx context.Context, f *Feature) (*Feature, error) {
	if f.Environments == nil {
		f.Environments = map[string]FeatureEnvironmentConfig{}
	}
	if f.Tags == nil {
		f.Tags = []string{}
	}
	if f.Prerequisites == nil {
		f.Prerequisites = []string{}
	}
	out, err := fetchOne[featureResponse](ctx, c, "POST", "/features", f, http.StatusOK, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	return &out.Feature, nil
}

// GetFeature fetches a feature by its ID.
func (c *Client) GetFeature(ctx context.Context, id string) (*Feature, error) {
	out, err := fetchOne[featureResponse](ctx, c, "GET", "/features/"+id, nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &out.Feature, nil
}

// UpdateFeature updates an existing feature by its ID.
func (c *Client) UpdateFeature(ctx context.Context, id string, f *Feature) (*Feature, error) {
	if f.Environments == nil {
		f.Environments = map[string]FeatureEnvironmentConfig{}
	}
	if f.Tags == nil {
		f.Tags = []string{}
	}
	if f.Prerequisites == nil {
		f.Prerequisites = []string{}
	}
	out, err := fetchOne[featureResponse](ctx, c, "POST", "/features/"+id, f, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &out.Feature, nil
}

// DeleteFeature removes a feature by its ID.
func (c *Client) DeleteFeature(ctx context.Context, id string) error {
	return c.remove(ctx, "/features/"+id, http.StatusOK, http.StatusNoContent)
}

// FindFeatureByName searches for a feature by its ID and returns the first match, handling pagination.
func (c *Client) FindFeatureByName(ctx context.Context, id string) (*Feature, error) {
	features, err := fetchAllPages[Feature](ctx, c, "GET", "/features", nil, "features", http.StatusOK)
	if err != nil {
		return nil, err
	}
	for _, f := range features {
		if f.ID == id {
			return &f, nil
		}
	}
	return nil, ErrNotFound
}
