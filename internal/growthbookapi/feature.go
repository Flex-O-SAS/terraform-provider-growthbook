package growthbookapi

import (
	"context"
)

// CreateFeature creates a new feature in GrowthBook.
func (c *Client) CreateFeature(ctx context.Context, f *Feature) (*Feature, error) {
	if f.Tags == nil {
		f.Tags = []string{}
	}
	if f.Prerequisites == nil {
		f.Prerequisites = []string{}
	}
	out, err := fetcher[Feature](c, "POST", "/features").One(ctx, f, "feature")
	if err != nil {
		return nil, err
	}
	if featureNeedsFollowUpUpdate(f) {
		id := out.ID
		if id == "" {
			id = f.ID
		}
		updated, err := c.UpdateFeature(ctx, id, f)
		if err != nil {
			return nil, err
		}
		return updated, nil
	}
	return &out, nil
}

// GetFeature fetches a feature by its ID.
func (c *Client) GetFeature(ctx context.Context, id string) (*Feature, error) {
	out, err := fetcher[Feature](c, "GET", "/features/"+id).One(ctx, nil, "feature")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateFeature updates an existing feature by its ID.
func (c *Client) UpdateFeature(ctx context.Context, id string, f *Feature) (*Feature, error) {
	if f.Tags == nil {
		f.Tags = []string{}
	}
	if f.Prerequisites == nil {
		f.Prerequisites = []string{}
	}
	out, err := fetcher[Feature](c, "POST", "/features/"+id).One(ctx, f, "feature")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteFeature removes a feature by its ID.
func (c *Client) DeleteFeature(ctx context.Context, id string) error {
	return c.delete(ctx, "/features/"+id)
}

// FindFeatureByName searches for a feature by its ID and returns the first match, handling pagination.
func (c *Client) FindFeatureByName(ctx context.Context, id string) (*Feature, error) {
	features, err := fetcher[Feature](c, "GET", "/features").All(ctx, nil, "features")
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

func featureNeedsFollowUpUpdate(f *Feature) bool {
	if f == nil {
		return false
	}
	for _, env := range f.Environments {
		for _, rule := range env.Rules {
			if len(rule.SavedGroupTargeting) > 0 || len(rule.ScheduleRules) > 0 {
				return true
			}
		}
	}
	return false
}
