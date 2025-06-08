package growthbookapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) CreateEnvironment(e *Environment) (*Environment, error) {
	resp, err := c.doRequest("POST", "/environments", e)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create environment failed: %s", string(b))
	}
	var out struct {
		Environment Environment `json:"environment"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out.Environment, nil
}

func (c *Client) UpdateEnvironment(id string, e *Environment) (*Environment, error) {
	resp, err := c.doRequest("PUT", "/environments/"+id, e)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("update environment failed: %s", string(b))
	}
	var out struct {
		Environment Environment `json:"environment"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out.Environment, nil
}

func (c *Client) DeleteEnvironment(id string) error {
	resp, err := c.doRequest("DELETE", "/environments/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete environment failed: %s", string(b))
	}
	return nil
}

// FindEnvironmentByID fetches an environment by its ID by listing all and filtering.
func (c *Client) FindEnvironmentByID(id string) (*Environment, error) {
	resp, err := c.doRequest("GET", "/environments", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list environments failed: %s", string(b))
	}
	var envs struct {
		Environments []Environment `json:"environments"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envs); err != nil {
		return nil, err
	}
	for _, env := range envs.Environments {
		if env.ID == id {
			return &env, nil
		}
	}
	return nil, ErrNotFound
}
