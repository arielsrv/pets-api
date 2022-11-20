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

func (c *Client) GetProject(projectID int64) (*responses.ProjectResponse, error) {
	response := c.rb.Get(fmt.Sprintf("/projects/%d", projectID))
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusOK {
		return nil, server.NewError(response.StatusCode, response.String())
	}

	var projectResponse responses.ProjectResponse
	err := response.FillUp(&projectResponse)
	if err != nil {
		return nil, err
	}

	return &projectResponse, nil
}

func NewGitLabClient(rb *rest.RequestBuilder) *Client {
	return &Client{rb: rb}
}

func (c *Client) GetGroups() ([]responses.GroupResponse, error) {
	response := c.rb.Get("/groups")
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusOK {
		return nil, server.NewError(response.StatusCode, response.String())
	}
	var groups []responses.GroupResponse
	err := response.FillUp(&groups)
	if err != nil {
		return nil, err
	}

	total, err := strconv.Atoi(response.Response.Header.Get("x-total-pages"))
	if err != nil {
		return nil, err
	}
	if total > 1 {
		var pages []*rest.FutureResponse
		c.rb.ForkJoin(func(c *rest.Concurrent) {
			for i := 2; i <= total; i++ {
				pages = append(pages, c.Get(fmt.Sprintf("/groups?page=%d", i)))
			}
		})
		for i := range pages {
			if pages[i].Response().StatusCode != http.StatusOK {
				return nil, server.NewError(response.StatusCode, response.String())
			}
			var page []responses.GroupResponse
			err = pages[i].Response().FillUp(&page)
			if err != nil {
				return nil, err
			}
			groups = append(groups, page...)
		}
	}

	return groups, nil
}

func (c *Client) CreateProject(request *requests.CreateProjectRequest) (*responses.CreateProjectResponse, error) {
	response := c.rb.Post("/projects", request)
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusCreated {
		return nil, server.NewError(response.StatusCode, response.String())
	}

	var createProjectResponse responses.CreateProjectResponse
	err := response.FillUp(&createProjectResponse)
	if err != nil {
		return nil, err
	}

	return &createProjectResponse, nil
}
