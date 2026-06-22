package client

import (
	"testing"

	"go.platform-mesh.io/golang-commons/directive/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewOpenFGAClient(t *testing.T) {
	client, err := NewOpenFGAClient(&mocks.OpenFGAServiceClient{})
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
