package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPostgres(t *testing.T) {
	// Initialize test config
	config := &SConfig{
		Host:         "localhost",
		Port:         "5432",
		User:         "admin",
		Password:     "admin",
		DatabaseName: "go_clean_architecture",
		SslMode:      "disable",
		TimeZone:     "Asia/Tehran",
	}

	// Test that NewPostgres initializes and returns an SPostgres object
	sPostgres := NewPostgres(config, nil)
	require.NotNil(t, sPostgres)
}

func TestSPostgres_Close(t *testing.T) {
	// Initialize test config
	config := &SConfig{
		Host:         "localhost",
		Port:         "5432",
		User:         "admin",
		Password:     "admin",
		DatabaseName: "go_clean_architecture",
		SslMode:      "disable",
		TimeZone:     "Asia/Tehran",
	}

	// Initialize SPostgres object
	sPostgres := NewPostgres(config, nil)
	assert.NotNil(t, sPostgres)

	// Test that Close method returns no error
	err := sPostgres.Close()
	assert.NoError(t, err)
}
