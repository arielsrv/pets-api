package clients

import (
	"fmt"
	"github.com/internal/clients/requests"
	"github.com/internal/clients/responses"
	"net/http"
	"strconv"

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
		page, headerErr := strconv.Atoi(response.Response.Header.Get("x-page"))
		if headerErr != nil {
			return nil, headerErr
		}
		for i := page + 1; i <= total; i++ {
			response = r.rb.Get(fmt.Sprintf("/groups?page=%d", i))
			if response.StatusCode != http.StatusOK {
				return nil, shared.NewError(response.StatusCode, response.String())
			}
			var pageGroup []responses.GroupResponse
			response.FillUp(&pageGroup)
			groups = append(groups, pageGroup...)
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
