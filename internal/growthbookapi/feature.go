package growthbookapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) CreateFeature(f *Feature) (*Feature, error) {
	resp, err := c.doRequest("POST", "/features", f)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create feature failed: %s", string(b))
	}
	var out struct {
		Feature Feature `json:"feature"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out.Feature, nil
}

func (c *Client) GetFeature(id string) (*Feature, error) {
	resp, err := c.doRequest("GET", "/features/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get feature failed: %s", string(b))
	}
	var out struct {
		Feature Feature `json:"feature"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out.Feature, nil
}

func (c *Client) UpdateFeature(id string, f *Feature) (*Feature, error) {
	resp, err := c.doRequest("PUT", "/features/"+id, f)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("update feature failed: %s", string(b))
	}
	var out struct {
		Feature Feature `json:"feature"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out.Feature, nil
}

func (c *Client) DeleteFeature(id string) error {
	resp, err := c.doRequest("DELETE", "/features/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete feature failed: %s", string(b))
	}
	return nil
}
