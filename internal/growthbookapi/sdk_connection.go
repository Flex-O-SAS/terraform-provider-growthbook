package growthbookapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) CreateSDKConnection(s *SDKConnection) (*SDKConnection, error) {
	resp, err := c.doRequest("POST", "/sdk-connections", s)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create sdkconnection failed: %s", string(b))
	}
	var out struct {
		SDKConnection SDKConnection `json:"sdkConnection"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.SDKConnection.Languages) != 0 {
		out.SDKConnection.Language = out.SDKConnection.Languages[0]
	}
	return &out.SDKConnection, nil
}

func (c *Client) GetSDKConnection(id string) (*SDKConnection, error) {
	resp, err := c.doRequest("GET", "/sdk-connections/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get sdkconnection failed: %s", string(b))
	}
	var out struct {
		SDKConnection SDKConnection `json:"sdkConnection"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.SDKConnection.Languages) != 0 {
		out.SDKConnection.Language = out.SDKConnection.Languages[0]
	}
	return &out.SDKConnection, nil
}

func (c *Client) UpdateSDKConnection(id string, s *SDKConnection) (*SDKConnection, error) {
	resp, err := c.doRequest("PUT", "/sdk-connections/"+id, s)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("update sdkconnection failed: %s", string(b))
	}
	var out struct {
		SDKConnection SDKConnection `json:"sdkConnection"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.SDKConnection.Languages) != 0 {
		out.SDKConnection.Language = out.SDKConnection.Languages[0]
	}
	return &out.SDKConnection, nil
}

func (c *Client) DeleteSDKConnection(id string) error {
	resp, err := c.doRequest("DELETE", "/sdk-connections/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete sdkconnection failed: %s", string(b))
	}
	return nil
}

// FindSDKConnectionByName searches for an SDK connection by its name and returns the first match.
func (c *Client) FindSDKConnectionByName(name string) (*SDKConnection, error) {
	resp, err := c.doRequest("GET", "/sdk-connections", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list sdk connections failed: %s", string(b))
	}
	var sdks struct {
		Connections []SDKConnection `json:"connections"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&sdks); err != nil {
		return nil, err
	}
	for _, s := range sdks.Connections {
		if s.Name == name {
			return &s, nil
		}
	}
	return nil, ErrNotFound // not found
}

// FindSDKConnectionByID fetches an SDK connection by its ID.
func (c *Client) FindSDKConnectionByID(id string) (*SDKConnection, error) {
	return c.GetSDKConnection(id)
}
