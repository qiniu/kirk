package kirksdk

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	KeyStateEnabled  = "enabled"
	KeyStateDisabled = "disabled"
)

const (
	DefaultAppHost   = "https://app-api.qiniu.com"
	DefaultIndexHost = "https://index.qiniu.com"
)

const (
	ProductQcos         = "qcos"
	ProductQcosVpnProxy = "qcos_vpnproxy"
	ProductQcosGates    = "qcos_gates"
)

type AppClient interface {
	GetAppInfo(ctx context.Context) (ret AppInfo, err error)
	CreateApp(ctx context.Context, appName string, args CreateAppArgs) (ret CreateAppRet, err error)
	DeleteApp(ctx context.Context, appURI string) (err error)
	ListSubApps(ctx context.Context) (ret []AppInfo, err error)
	ListManagedApps(ctx context.Context) (ret []AppInfo, err error)
	GetAppKeys(ctx context.Context, appURI string) (ret []KeyPair, err error)

	GetRegion(ctx context.Context, regionName string) (ret RegionInfo, err error)
	ListRegions(ctx context.Context) (ret []RegionInfo, err error)

	ListAlertMethods(ctx context.Context) (ret []AlertMethodInfo, err error)
	GetAlertMethod(ctx context.Context, id string) (ret AlertMethodInfo, err error)
	CreateAlertMethod(ctx context.Context, args CreateAlertMethodArgs) (ret AlertMethodInfo, err error)
	UpdateAlertMethod(ctx context.Context, id string, args UpdateAlertMethodArgs) (ret AlertMethodInfo, err error)
	DeleteAlertMethod(ctx context.Context, id string) (err error)

	GetIndexClient(ctx context.Context) (IndexClient, error)
	GetQcosClient(ctx context.Context, appURI string) (QcosClient, error)
}

type AppConfig struct {
	AccessKey string `json:"ak"`
	SecretKey string `json:"sk"`
	Host      string `json:"appd_host"`
	Logger    *logrus.Logger
	UserAgent string
	Transport http.RoundTripper
}

type CreateAppArgs struct {
	Title   string            `json:"title"`
	SpecURI string            `json:"specUri"`
	SpecVer uint32            `json:"specVer"`
	Region  string            `json:"region"`
	Vars    map[string]string `json:"vars"`
	Links   []string          `json:"links"`
}

type CreateAppRet struct {
	AppURI string `json:"appUri"`
}

type AppInfo struct {
	ID               uint32    `json:"id"`
	URI              string    `json:"uri"`
	Region           string    `json:"region"`
	Title            string    `json:"title"`
	Status           string    `json:"status"`
	RunMode          string    `json:"runMode,omitempty"`
	ParentURI        string    `json:"parentUri,omitempty"`
	CreationTime     time.Time `json:"ctime"`
	ModificationTime time.Time `json:"mtime"`
	AppExtendedInfo
}

type AppExtendedInfo struct {
	VendorURI string                 `json:"vendorUri,omitempty"`
	SpecURI   string                 `json:"specUri,omitempty"`
	SpecVer   uint32                 `json:"specVer,omitempty"`
	SetupMode string                 `json:"setupMode,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Vars      map[string]string      `json:"vars,omitempty"`
}

/*
	获得应用自身 access keys (access key id / secret key)
*/
type KeyPair struct {
	AccessKey string `json:"ak"`
	SecretKey string `json:"sk"`
	State     string `json:"state"`
}

type RegionInfo struct {
	Name     string            `json:"name"`
	Desc     string            `json:"desc"`
	Products map[string]string `json:"products"`
}

type AlertMethods struct {
	Methods []AlertMethodInfo `json:"methods"`
}

type AlertMethodInfo struct {
	ID          uint64 `json:"id"`
	Owner       uint32 `json:"owner"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	Nationality string `json:"nationality"`
	Code        string `json:"code"`
}

type CreateAlertMethodArgs struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	Nationality string `json:"nationality"`
	Code        string `json:"code"`
}

type UpdateAlertMethodArgs struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	Nationality string `json:"nationality"`
	Code        string `json:"code"`
}
