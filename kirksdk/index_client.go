package kirksdk

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"qiniupkg.com/x/rpc.v7"
)

type IndexConfig struct {
	AccessKey string
	SecretKey string
	Host      string
	RootApp   string
	AuthHost  string
	UserAgent string
	Transport http.RoundTripper
}

type indexClientImp struct {
	host   string
	client rpc.Client
}

func NewIndexClient(cfg IndexConfig) IndexClient {

	p := new(indexClientImp)
	cfg.Host = cleanHost(cfg.Host)
	p.host = cfg.Host

	cfg.Transport = newKirksdkTransport(cfg.UserAgent, cfg.Transport)

	p.client = rpc.Client{&http.Client{Transport: newAuthTokenTransport(cfg)}}

	return p
}

func (p *indexClientImp) ListRepo(ctx context.Context, username string) (repos []*Repo, err error) {
	err = p.client.Call(ctx, &repos, "GET", fmt.Sprintf("%s/api/%s/repos", p.host, username))
	return
}

func (p *indexClientImp) ListRepoTags(ctx context.Context, username, repo string) (tags []*Tag, err error) {
	err = p.client.Call(ctx, &tags, "GET", fmt.Sprintf("%s/api/%s/%s/tags", p.host, username, repo))
	return
}

func (p *indexClientImp) GetImageConfig(ctx context.Context, username, repo, reference string) (res *ImageConfig, err error) {
	err = p.client.Call(ctx, &res, "GET", fmt.Sprintf("%s/api/%s/%s/repo/%s", p.host, username, repo, reference))
	return
}

func (p *indexClientImp) DeleteRepoTag(ctx context.Context, username, repo, reference string) (err error) {
	err = p.client.Call(ctx, nil, "DELETE", fmt.Sprintf("%s/api/%s/%s/repo/%s", p.host, username, repo, reference))
	if httpCodeOf(err)/100 == 2 {
		err = nil
	}
	return
}

func (p *indexClientImp) CreateTagFromRepo(ctx context.Context, username, repo, tag string, from *ImageSpec) (result *ImageSpec, err error) {
	values := url.Values{
		"from":      {from.Username + "/" + from.Repo},
		"reference": {from.Reference},
	}
	url := fmt.Sprintf("%s/api/%s/%s/repo/%s?%s", p.host, username, repo, tag, values.Encode())
	err = p.client.CallWithJson(ctx, &result, "POST", url, nil)
	if httpCodeOf(err)/100 == 2 {
		err = nil
	}
	return
}

type httpCoder interface {
	HttpCode() int
}

func httpCodeOf(err error) int {
	if hc, ok := err.(httpCoder); ok {
		return hc.HttpCode()
	}
	return 0
}
