package services_test

import (
	"testing"

	"github.com/arielsrv/pets-api/src/main/app/services"
	"github.com/stretchr/testify/assert"
)

func TestPingService_Ping(t *testing.T) {
	pingService := services.NewPingService()
	actual := pingService.Ping()

	assert.Equal(t, "pong", actual)
}
