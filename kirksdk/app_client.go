package kirksdk

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"qiniupkg.com/kirk/kirksdk/mac"
	"qiniupkg.com/x/rpc.v7"
)

const appVersionPrefix = "/v1"

type appClientImp struct {
	accessKey string
	secretKey string
	host      string
	userAgent string
	client    rpc.Client
}

func NewAppClient(cfg AppConfig) AppClient {

	p := new(appClientImp)

	p.host = cleanHost(cfg.Host)

	p.userAgent = cfg.UserAgent
	cfg.Transport = newKirksdkTransport(cfg.UserAgent, cfg.Transport)

	if cfg.AccessKey == "" { // client used inside intranet
		p.client = rpc.Client{&http.Client{Transport: cfg.Transport}}
	} else {
		p.accessKey = cfg.AccessKey
		p.secretKey = cfg.SecretKey
		m := mac.New(cfg.AccessKey, cfg.SecretKey)
		p.client = rpc.Client{mac.NewClient(m, cfg.Transport)}
	}

	return p
}

/*
	获得应用基本信息（可公开）
*/
func (p *appClientImp) GetAppInfo(ctx context.Context) (ret AppInfo, err error) {
	url := fmt.Sprintf("%s%s/info", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

/*
	创建或同步应用
*/
func (p *appClientImp) CreateApp(ctx context.Context, appName string, args CreateAppArgs) (ret CreateAppRet, err error) {
	url := fmt.Sprintf("%s%s/apps/%s", p.host, appVersionPrefix, appName)
	err = p.client.CallWithJson(ctx, &ret, "POST", url, args)
	return
}

/*
	删除一个子应用
*/
func (p *appClientImp) DeleteApp(ctx context.Context, appURI string) (err error) {
	url := fmt.Sprintf("%s%s/apps/%s", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

/*
	列出调用者作为 parent 的所有子App
*/
func (p *appClientImp) ListSubApps(ctx context.Context) (ret []AppInfo, err error) {
	url := fmt.Sprintf("%s%s/children", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

/*
	列出调用者作为 Vendor 所管辖的 App
*/
func (p *appClientImp) ListManagedApps(ctx context.Context) (ret []AppInfo, err error) {
	url := fmt.Sprintf("%s%s/managed", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

/*
	获取 region 信息
*/
func (p *appClientImp) GetRegion(ctx context.Context, regionName string) (ret RegionInfo, err error) {
	url := fmt.Sprintf("%s%s/regions/%s", p.host, appVersionPrefix, regionName)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

/*
	列出所有 region
*/
func (p *appClientImp) ListRegions(ctx context.Context) (ret []RegionInfo, err error) {
	url := p.host + appVersionPrefix + "/regions"
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

/*
	获得其他应用 access keys / secret keys
*/
func (p *appClientImp) GetAppKeys(ctx context.Context, appURI string) (ret []KeyPair, err error) {
	url := fmt.Sprintf("%s%s/apps/%s/keys", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// GET /v1/alert/methods
func (p *appClientImp) ListAlertMethods(ctx context.Context) (ret []AlertMethodInfo, err error) {
	ret1 := AlertMethods{
		Methods: []AlertMethodInfo{},
	}

	url := fmt.Sprintf("%s%s/alert/methods", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret1, "GET", url)
	ret = ret1.Methods
	return
}

// GET /v1/alert/methods/<id>
func (p *appClientImp) GetAlertMethod(ctx context.Context, id string) (ret AlertMethodInfo, err error) {
	url := fmt.Sprintf("%s%s/alert/methods/%s", p.host, appVersionPrefix, id)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// POST /v1/alert/methods
func (p *appClientImp) CreateAlertMethod(ctx context.Context, args CreateAlertMethodArgs) (ret AlertMethodInfo, err error) {
	url := fmt.Sprintf("%s%s/alert/methods", p.host, appVersionPrefix)
	err = p.client.CallWithJson(ctx, &ret, "POST", url, args)
	return
}

// POST /v1/alert/methods/<id>
func (p *appClientImp) UpdateAlertMethod(ctx context.Context, id string, args UpdateAlertMethodArgs) (ret AlertMethodInfo, err error) {
	url := fmt.Sprintf("%s%s/alert/methods/%s", p.host, appVersionPrefix, id)
	err = p.client.CallWithJson(ctx, &ret, "POST", url, args)
	return
}

// DELETE /v1/alert/methods/<id>
func (p *appClientImp) DeleteAlertMethod(ctx context.Context, id string) (err error) {
	url := fmt.Sprintf("%s%s/alert/methods/%s", p.host, appVersionPrefix, id)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

func (p *appClientImp) GetIndexClient(ctx context.Context) (client IndexClient, err error) {
	appInfo, err := p.GetAppInfo(ctx)
	if err != nil {
		return
	}

	indexCfg := IndexConfig{
		AccessKey: p.accessKey,
		SecretKey: p.secretKey,
		Host:      DefaultIndexHost,
		RootApp:   appInfo.URI,
		UserAgent: p.userAgent,
	}

	return NewIndexClient(indexCfg), nil
}

func (p *appClientImp) GetQcosClient(ctx context.Context, appURI string) (client QcosClient, err error) {

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
		appInfos, err := p.ListSubApps(ctx)
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
