package kirksdk

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIndexGetConfig(t *testing.T) {
	config := IndexConfig{
		AccessKey: "ak",
		SecretKey: "sk",
		Host:      "https://index.test.url",
		RootApp:   "root",
		AuthHost:  "https://index.auth.test.url",
		UserAgent: "index.ua",
		Transport: http.DefaultTransport,
	}

	client := NewIndexClient(config)
	assert.EqualValues(t, config, client.GetConfig())
}
