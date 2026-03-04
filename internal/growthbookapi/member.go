package growthbookapi

import (
	"context"
)

// ListMembers fetches all organization members with pagination.
func (c *Client) ListMembers(ctx context.Context) ([]Member, error) {
	return fetcher[Member](c, "GET", "/members").All(ctx, nil, "members")
}

// GetMember fetches a member by ID (lists all and filters).
func (c *Client) GetMember(ctx context.Context, id string) (*Member, error) {
	members, err := c.ListMembers(ctx)
	if err != nil {
		return nil, err
	}
	for _, m := range members {
		if m.ID == id {
			return &m, nil
		}
	}
	return nil, ErrNotFound
}

// FindMemberByEmail searches for a member by email.
func (c *Client) FindMemberByEmail(ctx context.Context, email string) (*Member, error) {
	members, err := c.ListMembers(ctx)
	if err != nil {
		return nil, err
	}
	for _, m := range members {
		if m.Email == email {
			return &m, nil
		}
	}
	return nil, ErrNotFound
}

// UpdateMemberRole updates a member's global role, environments, and project roles.
func (c *Client) UpdateMemberRole(ctx context.Context, id string, update *MemberRoleUpdate) (*Member, error) {
	out, err := fetcher[Member](c, "POST", "/members/"+id+"/role").One(ctx, update, "member")
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMember removes a member from the organization.
func (c *Client) DeleteMember(ctx context.Context, id string) error {
	return c.delete(ctx, "/members/"+id)
}
