package gitlab

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/internal/clients/gitlab/requests"
	responses2 "github.com/internal/clients/gitlab/responses"

	"github.com/internal/shared"

	"github.com/arielsrv/golang-toolkit/rest"
)

type IGitLabClient interface {
	GetGroups() ([]responses2.GroupResponse, error)
	CreateProject(request *requests.CreateProjectRequest) (*responses2.CreateProjectResponse, error)
	GetProject(projectID int64) (*responses2.ProjectResponse, error)
}

type GitLabClient struct {
	rb *rest.RequestBuilder
}

func (r *GitLabClient) GetProject(projectID int64) (*responses2.ProjectResponse, error) {
	response := r.rb.Get(fmt.Sprintf("/projects/%d", projectID))
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusOK {
		return nil, shared.NewError(response.StatusCode, response.String())
	}

	var projectResponse responses2.ProjectResponse
	response.FillUp(&projectResponse)

	return &projectResponse, nil
}

func NewGitLabClient(rb *rest.RequestBuilder) *GitLabClient {
	return &GitLabClient{rb: rb}
}

func (r *GitLabClient) GetGroups() ([]responses2.GroupResponse, error) {
	response := r.rb.Get("/groups")
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusOK {
		return nil, shared.NewError(response.StatusCode, response.String())
	}
	var groups []responses2.GroupResponse
	response.FillUp(&groups)

	total, err := strconv.Atoi(response.Response.Header.Get("x-total-pages"))
	if err != nil {
		return nil, err
	}
	if total > 1 {
		var pages []*rest.FutureResponse
		r.rb.ForkJoin(func(c *rest.Concurrent) {
			for i := 2; i <= total; i++ {
				pages = append(pages, c.Get(fmt.Sprintf("/groups?page=%d", i)))
			}
		})
		for i := range pages {
			if pages[i].Response().StatusCode != http.StatusOK {
				return nil, shared.NewError(response.StatusCode, response.String())
			}
			var page []responses2.GroupResponse
			pages[i].Response().FillUp(&page)
			groups = append(groups, page...)
		}
	}

	return groups, nil
}

func (r *GitLabClient) CreateProject(request *requests.CreateProjectRequest) (*responses2.CreateProjectResponse, error) {
	response := r.rb.Post("/projects", request)
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusCreated {
		return nil, shared.NewError(response.StatusCode, response.String())
	}

	var createProjectResponse responses2.CreateProjectResponse
	response.FillUp(&createProjectResponse)

	return &createProjectResponse, nil
}
