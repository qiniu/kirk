package kirksdk

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAccountGetConfig(t *testing.T) {
	config := AccountConfig{
		AccessKey: "ak",
		SecretKey: "sk",
		Host:      "https://account.test.url",
		UserAgent: "account.ua",
		Transport: http.DefaultTransport,
	}

	client := NewAccountClient(config)
	assert.EqualValues(t, config, client.GetConfig())
}
