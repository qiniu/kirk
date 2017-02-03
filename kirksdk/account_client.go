package kirksdk

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"qiniupkg.com/kirk/kirksdk/mac"
	"qiniupkg.com/x/rpc.v7"
)

const appVersionPrefix = "/v3"

var ErrInvalidAppURI = errors.New("app uri is invalid")

type accountClientImp struct {
	config    AccountConfig
	accessKey string
	secretKey string
	host      string
	userAgent string
	client    rpc.Client
	transport http.RoundTripper
}

func NewAccountClient(cfg AccountConfig) AccountClient {

	p := new(accountClientImp)
	p.config = cfg
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

func (p *accountClientImp) GetConfig() (ret AccountConfig) {
	return p.config
}

func (p *accountClientImp) GetAccountInfo(ctx context.Context) (ret AccountInfo, err error) {
	url := fmt.Sprintf("%s%s/info", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

type createAppArgsWithName struct {
	Name string `json:"name"`
	CreateAppArgs
}

func (p *accountClientImp) CreateApp(ctx context.Context, appName string, args CreateAppArgs) (ret AppInfo, err error) {
	argsWithName := createAppArgsWithName{
		Name:          appName,
		CreateAppArgs: args,
	}
	url := fmt.Sprintf("%s%s/apps", p.host, appVersionPrefix)
	err = p.client.CallWithJson(ctx, &ret, "POST", url, argsWithName)
	return
}

func (p *accountClientImp) DeleteApp(ctx context.Context, appURI string) (err error) {
	url := fmt.Sprintf("%s%s/apps/%s", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

func (p *accountClientImp) GetApp(ctx context.Context, appURI string) (ret AppInfo, err error) {
	url := fmt.Sprintf("%s%s/apps/%s", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) GetAppKeys(ctx context.Context, appURI string) (ret []KeyPair, err error) {
	url := fmt.Sprintf("%s%s/apps/%s/keys", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) ListApps(ctx context.Context) (ret []AppInfo, err error) {
	url := fmt.Sprintf("%s%s/apps", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) GetAppQuota(ctx context.Context, appURI string) (ret []QuotaItem, err error) {
	quotaMap := make(map[string]string)
	url := fmt.Sprintf("%s%s/apps/%s/quota", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, &quotaMap, "GET", url)
	if err != nil {
		return
	}

	// parse quota items
	for k, v := range quotaMap {
		parts := strings.Split(v, "/")
		if len(parts) < 2 {
			continue
		}

		used, err := strconv.ParseInt(parts[0], 10, 0)
		if err != nil {
			return nil, ErrParseQuotaError
		}
		max, err := strconv.ParseInt(parts[1], 10, 0)
		if err != nil {
			return nil, ErrParseQuotaError
		}

		item := QuotaItem{
			Name: k,
			Used: used,
			Max:  max,
		}
		ret = append(ret, item)
	}

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
	url := fmt.Sprintf("%s%s/regions", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) CreateAlertMethod(ctx context.Context, appURI string, args CreateAlertMethodArgs) (ret AlertMethodInfo, err error) {
	url := fmt.Sprintf("%s%s/apps/%s/alert/methods", p.host, appVersionPrefix, appURI)
	err = p.client.CallWithJson(ctx, &ret, "POST", url, args)
	return
}

func (p *accountClientImp) DeleteAlertMethod(ctx context.Context, appURI string, id string) (err error) {
	url := fmt.Sprintf("%s%s/apps/%s/alert/methods/%s", p.host, appVersionPrefix, appURI, id)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

func (p *accountClientImp) GetAlertMethod(ctx context.Context, appURI string, id string) (ret AlertMethodInfo, err error) {
	url := fmt.Sprintf("%s%s/apps/%s/alert/methods/%s", p.host, appVersionPrefix, appURI, id)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) ListAlertMethod(ctx context.Context, appURI string) (ret []AlertMethodInfo, err error) {
	url := fmt.Sprintf("%s%s/apps/%s/alert/methods", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) UpdateAlertMethod(ctx context.Context, appURI string, id string, args UpdateAlertMethodArgs) (ret AlertMethodInfo, err error) {
	url := fmt.Sprintf("%s%s/apps/%s/alert/methods/%s", p.host, appVersionPrefix, appURI, id)
	err = p.client.CallWithJson(ctx, &ret, "PUT", url, args)
	return
}

func (p *accountClientImp) CreateAppGrant(ctx context.Context, appURI, username string) (err error) {
	url := fmt.Sprintf("%s%s/apps/%s/grants/%s", p.host, appVersionPrefix, appURI, username)
	err = p.client.Call(ctx, nil, "PUT", url)
	return
}

func (p *accountClientImp) ListGrants(ctx context.Context) (ret []GrantInfo, err error) {
	url := fmt.Sprintf("%s%s/grants", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) DeleteAppGrant(ctx context.Context, appURI, username string) (err error) {
	url := fmt.Sprintf("%s%s/apps/%s/grants/%s", p.host, appVersionPrefix, appURI, username)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

func (p *accountClientImp) ListAppGrantedUsers(ctx context.Context, appURI string) (ret []AppGrantedUser, err error) {
	url := fmt.Sprintf("%s%s/apps/%s/grants", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) ListGrantedApps(ctx context.Context) (ret []AppInfo, err error) {
	url := fmt.Sprintf("%s%s/granted", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) GetGrantedAppKey(ctx context.Context, appURI string) (ret GrantedAppKey, err error) {
	url := fmt.Sprintf("%s%s/granted/%s/key", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) GetAppspecs(ctx context.Context, specURI string) (ret SpecInfo, err error) {
	url := fmt.Sprintf("%s%s/appspecs/%s", p.host, appVersionPrefix, specURI)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) ListPublicspecs(ctx context.Context) (ret []SpecInfo, err error) {
	url := fmt.Sprintf("%s%s/publicspecs", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) ListGrantedspecs(ctx context.Context) (ret []SpecInfo, err error) {
	url := fmt.Sprintf("%s%s/grantedspecs", p.host, appVersionPrefix)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *accountClientImp) GetVendorManagedAppStatus(ctx context.Context, appURI string) (ret VendorManagedAppStatus, err error) {
	url := fmt.Sprintf("%s%s/apps/%s/status", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, &ret, "PUT", url)
	return
}

func (p *accountClientImp) GetVendorManagedAppEntry(ctx context.Context, appURI string) (ret VendorManagedAppEntry, err error) {
	url := fmt.Sprintf("%s%s/apps/%s/entry", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, &ret, "PUT", url)
	return
}

func (p *accountClientImp) VendorManagedAppRepair(ctx context.Context, appURI string) (err error) {
	url := fmt.Sprintf("%s%s/apps/%s/repair", p.host, appVersionPrefix, appURI)
	err = p.client.Call(ctx, nil, "PUT", url)
	return
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

	// app uri should follow the format: "username.appname"
	// or it will return an invalid app uri
	appURIParts := strings.Split(appURI, ".")
	if len(appURIParts) < 2 {
		return nil, ErrInvalidAppURI
	}

	// check if app is granted
	accountInfo, err := p.GetAccountInfo(ctx)
	if err != nil {
		return
	}
	isGranted := (accountInfo.Name != appURIParts[0])

	getAppKeyFunc := func() (result keyResult) {
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
		return
	}

	getGrantedAppKeyFunc := func() (result keyResult) {
		keyPair, err := p.GetGrantedAppKey(ctx, appURI)
		result.err = err
		if err == nil {
			result.ak = keyPair.Ak
			result.sk = keyPair.Sk
		}
		return
	}

	// set up list apps and get key func
	listAppsFunc := p.ListApps
	getKeyFunc := getAppKeyFunc
	if isGranted {
		listAppsFunc = p.ListGrantedApps
		getKeyFunc = getGrantedAppKeyFunc
	}

	keyChan := make(chan keyResult)
	endpointChan := make(chan endpointResult)

	// Get app access key & secret key
	go func() {
		result := getKeyFunc()

		if result.ak == "" && result.err == nil {
			result.err = fmt.Errorf("Fail to find keys for app \"%s\"", appURI)
		}

		keyChan <- result
	}()

	// Get qcos end point
	go func() {
		var result endpointResult
		appInfos, err := listAppsFunc(ctx)
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

		endpoint, ok := regionInfo.Products[ProductAPI]
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
