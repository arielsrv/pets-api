package gitlab

import (
	"fmt"
	"github.com/src/main/app/server"

	"github.com/src/main/app/clients/gitlab/responses"

	"net/http"
	"strconv"

	"github.com/arielsrv/golang-toolkit/rest"
	"github.com/src/main/app/clients/gitlab/requests"
)

type IGitLabClient interface {
	GetGroups() ([]responses.GroupResponse, error)
	CreateProject(request *requests.CreateProjectRequest) (*responses.CreateProjectResponse, error)
	GetProject(projectID int64) (*responses.ProjectResponse, error)
}

type Client struct {
	rb *rest.RequestBuilder
}

func (g *Client) GetProject(projectID int64) (*responses.ProjectResponse, error) {
	response := g.rb.Get(fmt.Sprintf("/projects/%d", projectID))
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusOK {
		return nil, server.NewError(response.StatusCode, response.String())
	}

	var projectResponse responses.ProjectResponse
	response.FillUp(&projectResponse)

	return &projectResponse, nil
}

func NewGitLabClient(rb *rest.RequestBuilder) *Client {
	return &Client{rb: rb}
}

func (g *Client) GetGroups() ([]responses.GroupResponse, error) {
	response := g.rb.Get("/groups")
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusOK {
		return nil, server.NewError(response.StatusCode, response.String())
	}
	var groups []responses.GroupResponse
	response.FillUp(&groups)

	total, err := strconv.Atoi(response.Response.Header.Get("x-total-pages"))
	if err != nil {
		return nil, err
	}
	if total > 1 {
		var pages []*rest.FutureResponse
		g.rb.ForkJoin(func(c *rest.Concurrent) {
			for i := 2; i <= total; i++ {
				pages = append(pages, c.Get(fmt.Sprintf("/groups?page=%d", i)))
			}
		})
		for i := range pages {
			if pages[i].Response().StatusCode != http.StatusOK {
				return nil, server.NewError(response.StatusCode, response.String())
			}
			var page []responses.GroupResponse
			pages[i].Response().FillUp(&page)
			groups = append(groups, page...)
		}
	}

	return groups, nil
}

func (g *Client) CreateProject(request *requests.CreateProjectRequest) (*responses.CreateProjectResponse, error) {
	response := g.rb.Post("/projects", request)
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusCreated {
		return nil, server.NewError(response.StatusCode, response.String())
	}

	var createProjectResponse responses.CreateProjectResponse
	response.FillUp(&createProjectResponse)

	return &createProjectResponse, nil
}
