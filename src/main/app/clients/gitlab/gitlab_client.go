package gitlab

import (
	"fmt"

	"github.com/src/main/app/config"
	"github.com/src/main/app/infrastructure/secrets"
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
	rb          *rest.RequestBuilder
	baseURL     string
	secretStore secrets.ISecretStore
}

func (c *Client) GetProject(projectID int64) (*responses.ProjectResponse, error) {
	err := addHeaders(c)
	if err != nil {
		return nil, err
	}
	apiURL := fmt.Sprintf("%s/projects/%d", c.baseURL, projectID)
	response := c.rb.Get(apiURL)
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusOK {
		return nil, server.NewError(response.StatusCode, response.String())
	}

	var projectResponse responses.ProjectResponse
	err = response.FillUp(&projectResponse)
	if err != nil {
		return nil, err
	}

	return &projectResponse, nil
}

func NewGitLabClient(rb *rest.RequestBuilder, secretStore secrets.ISecretStore) *Client {
	return &Client{
		baseURL:     config.String("rest.client.gitlab.baseURL"),
		rb:          rb,
		secretStore: secretStore,
	}
}

func (c *Client) GetGroups() ([]responses.GroupResponse, error) {
	apiURL := fmt.Sprintf("%s/groups", c.baseURL)
	err := addHeaders(c)
	if err != nil {
		return nil, err
	}
	response := c.rb.Get(apiURL)
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusOK {
		return nil, server.NewError(response.StatusCode, response.String())
	}
	var groups []responses.GroupResponse
	err = response.FillUp(&groups)
	if err != nil {
		return nil, err
	}

	total, err := strconv.Atoi(response.Response.Header.Get("x-total-pages"))
	if err != nil {
		return nil, err
	}
	if total > 1 {
		var pages []*rest.FutureResponse
		c.rb.ForkJoin(func(concurrent *rest.Concurrent) {
			for i := 2; i <= total; i++ {
				pageURL := fmt.Sprintf("%s/groups?page=%d", c.baseURL, i)
				pages = append(pages, concurrent.Get(pageURL))
			}
		})
		for i := range pages {
			if pages[i].Response().Err != nil {
				return nil, pages[i].Response().Err
			}
			if pages[i].Response().StatusCode != http.StatusOK {
				return nil, server.NewError(pages[i].Response().StatusCode, pages[i].Response().String())
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
	err := addHeaders(c)
	if err != nil {
		return nil, err
	}
	apiURL := fmt.Sprintf("%s/projects", c.baseURL)
	response := c.rb.Post(apiURL, request)
	if response.Err != nil {
		return nil, response.Err
	}
	if response.StatusCode != http.StatusCreated {
		return nil, server.NewError(response.StatusCode, response.String())
	}

	var createProjectResponse responses.CreateProjectResponse
	err = response.FillUp(&createProjectResponse)
	if err != nil {
		return nil, err
	}

	return &createProjectResponse, nil
}

func addHeaders(c *Client) error {
	secret := c.secretStore.GetSecret("SECRETS_STORE_PETS-API_GITLAB_TOKEN_KEY_NAME")
	if secret.Err != nil {
		return secret.Err
	}
	headers := http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", secret.Value)},
	}
	c.rb.Headers = headers

	return nil
}
