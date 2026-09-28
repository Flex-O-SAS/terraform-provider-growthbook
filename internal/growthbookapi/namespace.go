package growthbookapi

import "context"

func (c *Client) CreateNamespace(ctx context.Context, ns *Namespace) (*Namespace, error) {
	out, err := fetcher[Namespace](c, "POST", "/namespaces").One(ctx, ns, "namespace")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetNamespace(ctx context.Context, id string) (*Namespace, error) {
	out, err := fetcher[Namespace](c, "GET", "/namespaces/"+id).One(ctx, nil, "namespace")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateNamespace(ctx context.Context, id string, ns *Namespace) (*Namespace, error) {
	out, err := fetcher[Namespace](c, "PUT", "/namespaces/"+id).One(ctx, ns, "namespace")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteNamespace(ctx context.Context, id string) error {
	return c.delete(ctx, "/namespaces/"+id)
}

func (c *Client) FindNamespaceByDisplayName(ctx context.Context, displayName string) (*Namespace, error) {
	items, err := fetcher[Namespace](c, "GET", "/namespaces").All(ctx, nil, "namespaces")
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.DisplayName == displayName {
			return &item, nil
		}
	}
	return nil, ErrNotFound
}
