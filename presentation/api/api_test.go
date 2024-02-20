package api

import (
	"github.com/ehsandavari/go-context-plus"
	jwtMocks "github.com/ehsandavari/go-jwt/mocks"
	"github.com/ehsandavari/go-logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"health-check/infrastructure/config"
	"io"
	"net/http"
	"testing"
)

func Test_NewSApi_Start_Stop(t *testing.T) {
	isEnabled := true
	sConfig := &config.SConfig{
		Service: &config.SService{
			Api: &config.SApi{
				IsEnabled: &isEnabled,
				Host:      "localhost",
				Port:      "8080",
				Mode:      "debug",
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewMockILogger(ctrl)
	mockJwt := jwtMocks.NewMockIJwtServer(ctrl)

	sApi := NewSApi(nil, sConfig, mockJwt, mockLogger, nil)

	mockLogger.EXPECT().WithAny("api server info", sConfig.Service.Api).Return(mockLogger).Times(1)
	mockLogger.EXPECT().Info(contextplus.Background(), "api server start").Times(1)

	// Start the server
	sApi.Start()

	// Act
	resp, err := http.Get("http://localhost:8080/-/health")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			assert.NoError(t, err)
		}
	}(resp.Body)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Stop the server
	sApi.Stop()

	_, err = http.Get("http://localhost:8080/-/health")
	assert.Error(t, err)
}
