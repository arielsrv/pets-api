package clients

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/internal/clients/requests"
	"github.com/internal/clients/responses"

	"github.com/internal/shared"

	"github.com/arielsrv/golang-toolkit/rest"
)

type IGitLabClient interface {
	GetGroups() ([]responses.GroupResponse, error)
	CreateProject(request *requests.CreateProjectRequest) error
}

type GitLabClient struct {
	rb *rest.RequestBuilder
}

func NewGitLabClient(rb *rest.RequestBuilder) *GitLabClient {
	return &GitLabClient{rb: rb}
}

func (r *GitLabClient) GetGroups() ([]responses.GroupResponse, error) {
	response := r.rb.Get("/groups")
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusOK {
		return nil, shared.NewError(response.StatusCode, response.String())
	}
	var groups []responses.GroupResponse
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
			var page []responses.GroupResponse
			pages[i].Response().FillUp(&page)
			groups = append(groups, page...)
		}
	}

	return groups, nil
}

func (r *GitLabClient) CreateProject(request *requests.CreateProjectRequest) error {
	response := r.rb.Post("/projects", request)
	if response.Err != nil {
		return response.Err
	}
	if response.StatusCode != http.StatusCreated {
		return shared.NewError(response.StatusCode, response.String())
	}
	return nil
}
