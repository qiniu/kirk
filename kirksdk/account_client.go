package kirksdk

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/context"
	"qiniupkg.com/kirk/kirksdk/mac"
	"qiniupkg.com/x/rpc.v7"
)

const appVersionPrefix = "/v1"

type appClient struct {
	appURI    string
	accessKey string
	secretKey string
	host      string
	userAgent string
	client    rpc.Client
}

func (p *appClient) getInfo(ctx context.Context) (ret AppInfo, err error) {
	url := fmt.Sprintf("%s%s/info", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *appClient) listAlertMethods(ctx context.Context) (ret []AlertMethodInfo, err error) {
	ret1 := AlertMethods{
		Methods: []AlertMethodInfo{},
	}

	url := fmt.Sprintf("%s%s/alert/methods", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret1, "GET", url)
	ret = ret1.Methods
	return
}

func (p *appClient) getAlertMethod(ctx context.Context, id string) (ret AlertMethodInfo, err error) {
	url := fmt.Sprintf("%s%s/alert/methods/%s", p.host, appVersionPrefix, id)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *appClient) createAlertMethod(ctx context.Context, args CreateAlertMethodArgs) (ret AlertMethodInfo, err error) {
	url := fmt.Sprintf("%s%s/alert/methods", p.host, appVersionPrefix)
	err = p.client.CallWithJson(ctx, &ret, "POST", url, args)
	return
}

func (p *appClient) updateAlertMethod(ctx context.Context, id string, args UpdateAlertMethodArgs) (ret AlertMethodInfo, err error) {
	url := fmt.Sprintf("%s%s/alert/methods/%s", p.host, appVersionPrefix, id)
	err = p.client.CallWithJson(ctx, &ret, "POST", url, args)
	return
}

func (p *appClient) deleteAlertMethod(ctx context.Context, id string) (err error) {
	url := fmt.Sprintf("%s%s/alert/methods/%s", p.host, appVersionPrefix, id)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

type accountClientImp struct {
	accessKey string
	secretKey string
	host      string
	userAgent string
	client    rpc.Client

	transport     http.RoundTripper
	appsClientMap map[string]*appClient
	mapLock       *sync.Mutex
}

func NewAccountClient(cfg AccountConfig) AccountClient {

	p := new(accountClientImp)
	p.appsClientMap = make(map[string]*appClient)
	p.mapLock = &sync.Mutex{}
	p.host = cleanHost(cfg.Host)
	p.transport = cfg.Transport
	p.userAgent = cfg.UserAgent

	cfg.Transport = newKirksdkTransport(cfg.UserAgent, cfg.Transport)
	if cfg.AccessKey == "" {
		p.client = rpc.Client{&http.Client{Transport: cfg.Transport}}
	} else {
		p.accessKey = cfg.AccessKey
		p.secretKey = cfg.SecretKey
		m := mac.New(cfg.AccessKey, cfg.SecretKey)
		p.client = rpc.Client{mac.NewClient(m, cfg.Transport)}
	}

	return p
}

func (p *accountClientImp) getAppClient(ctx context.Context, appURI string) (ret *appClient, err error) {
	p.mapLock.Lock()
	defer p.mapLock.Unlock()

	var ok bool
	if ret, ok = p.appsClientMap[appURI]; !ok {
		ret, err = p.createAppClient(ctx, appURI)
		if err != nil {
			return
		}

		p.appsClientMap[appURI] = ret
	}

	return
}

func (p *accountClientImp) createAppClient(ctx context.Context, appURI string) (ret *appClient, err error) {
	keyPairs, err := p.GetAppKeys(ctx, appURI)
	if err != nil {
		return
	}

	var ak, sk string
	for _, keyPair := range keyPairs {
		if keyPair.State == KeyStateEnabled {
			ak = keyPair.AccessKey
			sk = keyPair.SecretKey
			break
		}
	}

	if ak == "" {
		err = fmt.Errorf("Fail to find keys for app \"%s\"", appURI)
		return
	}

	t := newKirksdkTransport(p.userAgent, p.transport)
	m := mac.New(ak, sk)
	c := rpc.Client{mac.NewClient(m, t)}

	return &appClient{
		appURI:    appURI,
		accessKey: ak,
		secretKey: sk,
		host:      p.host,
		userAgent: p.userAgent,
		client:    c,
	}, nil
}

func (p *accountClientImp) GetAccountInfo(ctx context.Context) (ret AccountInfo, err error) {
	url := fmt.Sprintf("%s%s/info", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

type createAppRet struct {
	AppURI string `json:"appUri"`
}

func (p *accountClientImp) CreateApp(ctx context.Context, appName string, args CreateAppArgs) (ret AppInfo, err error) {
	var createdURI createAppRet
	url := fmt.Sprintf("%s%s/apps/%s", p.host, appVersionPrefix, appName)
	err = p.client.CallWithJson(ctx, &createdURI, "POST", url, args)
	if err != nil {
		return
	}

	client, err := p.getAppClient(ctx, createdURI.AppURI)
	if err != nil {
		return
	}

	ret, err = client.getInfo(ctx)
	return
}

func (p *accountClientImp) DeleteApp(ctx context.Context, appURI string) (err error) {
	url := fmt.Sprintf("%s%s/apps/%s", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

func (p *accountClientImp) GetApp(ctx context.Context, appURI string) (ret AppInfo, err error) {
	c, err := p.getAppClient(ctx, appURI)
	if err != nil {
		return
	}

	return c.getInfo(ctx)
}

func (p *accountClientImp) GetAppKeys(ctx context.Context, appURI string) (ret []KeyPair, err error) {
	url := fmt.Sprintf("%s%s/apps/%s/keys", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) ListApps(ctx context.Context) (ret []AppInfo, err error) {
	url := fmt.Sprintf("%s%s/children", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) ListManagedApps(ctx context.Context) (ret []AppInfo, err error) {
	url := fmt.Sprintf("%s%s/managed", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) GetRegion(ctx context.Context, regionName string) (ret RegionInfo, err error) {
	url := fmt.Sprintf("%s%s/regions/%s", p.host, appVersionPrefix, regionName)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) ListRegions(ctx context.Context) (ret []RegionInfo, err error) {
	url := p.host + appVersionPrefix + "/regions"
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) CreateAlertMethod(ctx context.Context, appURI string, args CreateAlertMethodArgs) (ret AlertMethodInfo, err error) {
	c, err := p.getAppClient(ctx, appURI)
	if err != nil {
		return
	}

	return c.createAlertMethod(ctx, args)
}

func (p *accountClientImp) DeleteAlertMethod(ctx context.Context, appURI string, id string) (err error) {
	c, err := p.getAppClient(ctx, appURI)
	if err != nil {
		return
	}

	return c.deleteAlertMethod(ctx, id)
}

func (p *accountClientImp) GetAlertMethod(ctx context.Context, appURI string, id string) (ret AlertMethodInfo, err error) {
	c, err := p.getAppClient(ctx, appURI)
	if err != nil {
		return
	}

	return c.getAlertMethod(ctx, id)
}

func (p *accountClientImp) ListAlertMethod(ctx context.Context, appURI string) (ret []AlertMethodInfo, err error) {
	c, err := p.getAppClient(ctx, appURI)
	if err != nil {
		return
	}

	return c.listAlertMethods(ctx)
}

func (p *accountClientImp) UpdateAlertMethod(ctx context.Context, appURI string, id string, args UpdateAlertMethodArgs) (ret AlertMethodInfo, err error) {
	c, err := p.getAppClient(ctx, appURI)
	if err != nil {
		return
	}

	return c.updateAlertMethod(ctx, id, args)
}

func (p *accountClientImp) GetIndexClient(ctx context.Context) (client IndexClient, err error) {
	accountInfo, err := p.GetAccountInfo(ctx)
	if err != nil {
		return
	}

	indexCfg := IndexConfig{
		AccessKey: p.accessKey,
		SecretKey: p.secretKey,
		Host:      DefaultIndexHost,
		RootApp:   accountInfo.Name,
		UserAgent: p.userAgent,
	}

	return NewIndexClient(indexCfg), nil
}

func (p *accountClientImp) GetQcosClient(ctx context.Context, appURI string) (client QcosClient, err error) {

	type keyResult struct {
		ak  string
		sk  string
		err error
	}

	type endpointResult struct {
		endpoint string
		err      error
	}

	keyChan := make(chan keyResult)
	endpointChan := make(chan endpointResult)

	// Get app access key & secret key
	go func() {
		var result keyResult
		keyPairs, err := p.GetAppKeys(ctx, appURI)
		if err != nil {
			result.err = err
		} else {
			// Find an enabled KeyPairs
			for _, keyPair := range keyPairs {
				if keyPair.State == KeyStateEnabled {
					result.ak = keyPair.AccessKey
					result.sk = keyPair.SecretKey
					break
				}
			}
		}

		if result.ak == "" {
			result.err = fmt.Errorf("Fail to find keys for app \"%s\"", appURI)
		}

		keyChan <- result
	}()

	// Get qocos end point
	go func() {
		var result endpointResult
		appInfos, err := p.ListApps(ctx)
		if err != nil {
			result.err = err
			endpointChan <- result
			return
		}

		var region string
		// find the app
		for _, appInfo := range appInfos {
			if appInfo.URI == appURI {
				region = appInfo.Region
				break
			}
		}

		if region == "" {
			result.err = fmt.Errorf("Fail to find sub-app \"%s\"", appURI)
			endpointChan <- result
			return
		}

		regionInfo, err := p.GetRegion(context.TODO(), region)
		if err != nil {
			result.err = err
			endpointChan <- result
			return
		}

		endpoint, ok := regionInfo.Products[ProductQcos]
		if !ok {
			result.err = fmt.Errorf("Fail to find qcos endpoint of app \"%s\"", appURI)
		} else {
			result.endpoint = endpoint
		}

		endpointChan <- result
		return
	}()

	kr := <-keyChan
	if kr.err != nil {
		err = kr.err
		return
	}

	er := <-endpointChan
	if er.err != nil {
		err = er.err
		return
	}

	qcosCfg := QcosConfig{
		AccessKey: kr.ak,
		SecretKey: kr.sk,
		Host:      er.endpoint,
		UserAgent: p.userAgent,
	}

	return NewQcosClient(qcosCfg), nil
}
