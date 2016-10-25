package kirksdk

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"qiniupkg.com/kirk/kirksdk/mac"
	"qiniupkg.com/x/rpc.v7"
)

type IndexAuthConfig struct {
	AccessKey string
	SecretKey string
	Host      string
	UserAgent string
	Transport http.RoundTripper
}

type indexAuthClientImp struct {
	Host   string
	client rpc.Client
}

func NewIndexAuthClient(cfg IndexAuthConfig) IndexAuthClient {

	p := new(indexAuthClientImp)

	p.Host = cleanHost(cfg.Host)

	cfg.Transport = newKirksdkTransport(cfg.UserAgent, cfg.Transport)

	if cfg.AccessKey == "" { // client used inside intranet
		p.client = rpc.Client{&http.Client{Transport: cfg.Transport}}
	} else {
		m := mac.New(cfg.AccessKey, cfg.SecretKey)
		p.client = rpc.Client{mac.NewClient(m, cfg.Transport)}
	}

	return p
}

func (p *indexAuthClientImp) RequestAuthToken(ctx context.Context, scopes []string) (AuthToken, error) {
	param := url.Values{"scope": scopes}
	token := new(AuthToken)
	err := p.client.Call(ctx, token, "GET", fmt.Sprintf("%s/token?%s", p.Host, param.Encode()))

	return *token, err
}
