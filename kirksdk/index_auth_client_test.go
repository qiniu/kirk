package kirksdk

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIndexAuthGetConfig(t *testing.T) {
	config := IndexAuthConfig{
		AccessKey: "ak",
		SecretKey: "sk",
		Host:      "https://index.auth.test.url",
		UserAgent: "index.ua",
		Transport: http.DefaultTransport,
	}

	client := NewIndexAuthClient(config)
	assert.EqualValues(t, config, client.GetConfig())
}
