package clients

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/internal/model/requests"
	"github.com/internal/shared"

	"github.com/arielsrv/golang-toolkit/rest"
	"github.com/internal/model/responses"
)

type GitLabClient struct {
	rb *rest.RequestBuilder
}

func NewGitLabClient(rb *rest.RequestBuilder) *GitLabClient {
	return &GitLabClient{rb: rb}
}

func (r *GitLabClient) GetGroup(groupID int64) responses.GitLabGroupResponse {
	response := r.rb.Get(fmt.Sprintf("/groups/%d", groupID))
	if response.StatusCode != http.StatusOK {
		log.Fatal(response.StatusCode)
	}
	var group responses.GitLabGroupResponse
	response.FillUp(&group)

	return group
}

func (r *GitLabClient) GetGroups() ([]responses.GitLabGroupResponse, error) {
	response := r.rb.Get("/groups")
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusOK {
		return nil, shared.NewError(response.StatusCode, response.String())
	}
	var groups []responses.GitLabGroupResponse
	response.FillUp(&groups)

	total, err := strconv.Atoi(response.Response.Header.Get("x-total-pages"))
	if err != nil {
		log.Fatal(err)
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
			var pageGroup []responses.GitLabGroupResponse
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
