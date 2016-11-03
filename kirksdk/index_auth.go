package kirksdk

import (
	"errors"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

func newAuthTokenTransport(cfg IndexConfig) http.RoundTripper {

	if cfg.AuthHost == "" {
		cfg.AuthHost = cfg.Host
	}

	authCfg := IndexAuthConfig{
		AccessKey: cfg.AccessKey,
		SecretKey: cfg.SecretKey,
		Host:      cfg.AuthHost,
		Transport: cfg.Transport,
		UserAgent: cfg.UserAgent,
	}

	authClient := NewIndexAuthClient(authCfg)

	return &authTokenTransport{
		scopes: []string{
			"repository:" + cfg.RootApp + "/*:pull,push,del",
			"repository:library/*:pull",
		},
		lock:       make(chan struct{}, 1),
		token:      AuthToken{},
		authClient: authClient,
		transport:  cfg.Transport,
	}
}

type authTokenTransport struct {
	scopes     []string
	token      AuthToken
	lock       chan struct{}
	authClient IndexAuthClient
	transport  http.RoundTripper
}

func (p *authTokenTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	token := p.token
	if p.isTokenExpired(token.IssuedAt.Unix() + token.ExpiresIn) {
		tried := 0
	outFor:
		for {
			select {
			case p.lock <- struct{}{}:
				token = p.token
				if !p.isTokenExpired(token.IssuedAt.Unix() + token.ExpiresIn) {
					<-p.lock
					goto final
				}
				break outFor
			case <-time.After(1 * time.Second):
				tried++
				if tried > 3 {
					err = errors.New("Auth.Token, tried to refresh, but failed: times reached")
					return
				}
			}
		}

		token, err = p.authClient.RequestAuthToken(context.TODO(), p.scopes)
		if err != nil {
			time.Sleep(1 * time.Second)
			<-p.lock
			return
		}
		p.token = token
		<-p.lock
	}

final:
	req.Header.Set("Authorization", "Bearer "+token.Token)
	resp, err = p.transport.RoundTrip(req)
	return
}

func (p *authTokenTransport) isTokenExpired(tokenExpiry int64) bool {
	// refresh token in advance
	return time.Now().Unix() >= (tokenExpiry - 60)
}
