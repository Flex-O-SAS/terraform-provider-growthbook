package growthbookapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) CreateProject(p *Project) (*Project, error) {
	resp, err := c.doRequest("POST", "/projects", p)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create project failed: %s", string(b))
	}
	var out struct {
		Project Project `json:"project"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out.Project, nil
}

func (c *Client) GetProject(id string) (*Project, error) {
	resp, err := c.doRequest("GET", "/projects/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get project failed: %s", string(b))
	}
	var out struct {
		Project Project `json:"project"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out.Project, nil
}

func (c *Client) UpdateProject(id string, p *Project) (*Project, error) {
	resp, err := c.doRequest("PUT", "/projects/"+id, p)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("update project failed: %s", string(b))
	}
	var out struct {
		Project Project `json:"project"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out.Project, nil
}

func (c *Client) DeleteProject(id string) error {
	resp, err := c.doRequest("DELETE", "/projects/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete project failed: %s", string(b))
	}
	return nil
}

// FindProjectByName searches for a project by its name and returns the first match.
func (c *Client) FindProjectByName(name string) (*Project, error) {
	resp, err := c.doRequest("GET", "/projects", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list projects failed: %s", string(b))
	}
	var projects struct {
		Projects []Project `json:"projects"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}
	for _, p := range projects.Projects {
		if p.Name == name {
			return &p, nil
		}
	}
	return nil, ErrNotFound
}
