package growthbookapi

import "context"

func (c *Client) GetDataSource(ctx context.Context, id string) (*DataSource, error) {
	out, err := fetcher[DataSource](c, "GET", "/data-sources/"+id).One(ctx, nil, "dataSource")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) FindDataSourceByName(ctx context.Context, name string) (*DataSource, error) {
	items, err := fetcher[DataSource](c, "GET", "/data-sources").All(ctx, nil, "dataSources")
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.Name == name {
			return &item, nil
		}
	}
	return nil, ErrNotFound
}
