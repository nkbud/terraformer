// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rundeck

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
)

type JobGenerator struct {
	RundeckService
}

// RundeckJob represents a simplified Rundeck job structure
type RundeckJob struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Project     string `json:"project"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

// RundeckProject represents a Rundeck project
type RundeckProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (g *JobGenerator) InitResources() error {
	url := g.Args["url"].(string)
	token := g.Args["token"].(string)
	username := g.Args["username"].(string)
	password := g.Args["password"].(string)
	apiVersion := g.Args["api_version"].(string)

	if apiVersion == "" {
		apiVersion = "38" // Default API version
	}

	// First, get all projects
	projects, err := g.getProjects(url, token, username, password, apiVersion)
	if err != nil {
		return fmt.Errorf("failed to get projects: %v", err)
	}

	// For each project, get all jobs
	for _, project := range projects {
		jobs, err := g.getJobs(url, token, username, password, apiVersion, project.Name)
		if err != nil {
			if g.Verbose {
				fmt.Printf("Warning: failed to get jobs for project %s: %v\n", project.Name, err)
			}
			continue
		}

		g.Resources = append(g.Resources, g.createJobResources(jobs)...)
	}

	return nil
}

func (g *JobGenerator) getProjects(url, token, username, password, apiVersion string) ([]RundeckProject, error) {
	client := &http.Client{}
	
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/%s/projects", url, apiVersion), nil)
	if err != nil {
		return nil, err
	}

	// Set authentication
	if token != "" {
		req.Header.Set("X-Rundeck-Auth-Token", token)
	} else if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var projects []RundeckProject
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (g *JobGenerator) getJobs(url, token, username, password, apiVersion, project string) ([]RundeckJob, error) {
	client := &http.Client{}
	
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/%s/project/%s/jobs", url, apiVersion, project), nil)
	if err != nil {
		return nil, err
	}

	// Set authentication
	if token != "" {
		req.Header.Set("X-Rundeck-Auth-Token", token)
	} else if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jobs []RundeckJob
	err = json.Unmarshal(body, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (g *JobGenerator) createJobResources(jobs []RundeckJob) []terraformutils.Resource {
	var resources []terraformutils.Resource
	
	for _, job := range jobs {
		resourceName := g.normalizeResourceName(job.Project, job.Group, job.Name)
		
		resources = append(resources, terraformutils.NewSimpleResource(
			job.ID,
			resourceName,
			"rundeck_job",
			"rundeck",
			[]string{},
		))
	}
	
	return resources
}

func (g *JobGenerator) normalizeResourceName(project, group, name string) string {
	var parts []string
	
	if project != "" {
		parts = append(parts, g.sanitizeName(project))
	}
	
	if group != "" {
		parts = append(parts, g.sanitizeName(group))
	}
	
	if name != "" {
		parts = append(parts, g.sanitizeName(name))
	}
	
	return strings.Join(parts, "_")
}

func (g *JobGenerator) sanitizeName(name string) string {
	// Replace non-alphanumeric characters with underscores
	result := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, name)
	
	// Ensure it starts with a letter or underscore
	if len(result) > 0 && result[0] >= '0' && result[0] <= '9' {
		result = "_" + result
	}
	
	return result
}