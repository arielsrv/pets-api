package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arielsrv/pets-api/src/main/app/handlers"
	"github.com/arielsrv/pets-api/src/main/app/server"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PingHandlerSuite struct {
	suite.Suite
	app         *server.App
	pingHandler handlers.IPingHandler
	pingService *MockPingService
}

func (s *PingHandlerSuite) SetupTest() {
	s.pingService = new(MockPingService)
	s.pingHandler = handlers.NewPingHandler(s.pingService)
	s.app = server.New()
	s.app.Add(http.MethodGet, "/ping", s.pingHandler.Ping)
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(PingHandlerSuite))
}

type MockPingService struct {
	mock.Mock
}

func (mock *MockPingService) Ping() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func (s *PingHandlerSuite) TestPingHandler_Ping() {
	s.pingService.On("Ping").Return("pong")

	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	response, err := s.app.Test(request)
	s.Require().NoError(err)
	s.NotNil(response)
	s.Equal(http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	s.Require().NoError(err)
	s.NotNil(body)

	s.Equal("pong", string(body))
}
