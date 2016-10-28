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
	DefaultAccountHost = "https://app-api.qiniu.com"
	DefaultIndexHost   = "https://index.qiniu.com"
)

const (
	ProductQcos         = "qcos"
	ProductQcosVpnProxy = "qcos_vpnproxy"
	ProductQcosGates    = "qcos_gates"
)

type AccountClient interface {
	GetAccountInfo(ctx context.Context) (ret AccountInfo, err error)

	CreateApp(ctx context.Context, appName string, args CreateAppArgs) (ret AppInfo, err error)
	DeleteApp(ctx context.Context, appURI string) (err error)
	GetApp(ctx context.Context, appName string) (ret AppInfo, err error)
	GetAppKeys(ctx context.Context, appURI string) (ret []KeyPair, err error)
	ListApps(ctx context.Context) (ret []AppInfo, err error)
	ListManagedApps(ctx context.Context) (ret []AppInfo, err error)

	GetRegion(ctx context.Context, regionName string) (ret RegionInfo, err error)
	ListRegions(ctx context.Context) (ret []RegionInfo, err error)

	CreateAlertMethod(ctx context.Context, appURI string, args CreateAlertMethodArgs) (ret AlertMethodInfo, err error)
	DeleteAlertMethod(ctx context.Context, appURI string, id string) (err error)
	GetAlertMethod(ctx context.Context, appURI string, id string) (ret AlertMethodInfo, err error)
	ListAlertMethod(ctx context.Context, appURI string) (ret []AlertMethodInfo, err error)
	UpdateAlertMethod(ctx context.Context, appURI string, id string, args UpdateAlertMethodArgs) (ret AlertMethodInfo, err error)

	GetIndexClient(ctx context.Context) (client IndexClient, err error)
	GetQcosClient(ctx context.Context, appURI string) (client QcosClient, err error)
}

type AccountConfig struct {
	AccessKey string `json:"ak"`
	SecretKey string `json:"sk"`
	Host      string `json:"appd_host"`
	UserAgent string
	Logger    *logrus.Logger
	Transport http.RoundTripper
}

type CreateAppArgs struct {
	Title   string `json:"title"`
	Region  string `json:"region"`
	SpecURI string `json:"specUri"`
	SpecVer uint32 `json:"specVer"`
}

type AccountInfo struct {
	ID               uint32    `json:"id"`
	Name             string    `json:"uri"`
	Title            string    `json:"title"`
	CreationTime     time.Time `json:"ctime"`
	ModificationTime time.Time `json:"mtime"`
}

type AppInfo struct {
	ID               uint32    `json:"id"`
	URI              string    `json:"uri"`
	Region           string    `json:"region"`
	Title            string    `json:"title"`
	Account          string    `json:"parentUri,omitempty"`
	Status           string    `json:"status"`
	RunMode          string    `json:"runMode,omitempty"`
	CreationTime     time.Time `json:"ctime"`
	ModificationTime time.Time `json:"mtime"`
	AppExtendedInfo
}

type AppExtendedInfo struct {
	VendorURI string `json:"vendorUri,omitempty"`
	SpecURI   string `json:"specUri,omitempty"`
	SpecVer   uint32 `json:"specVer,omitempty"`
	SetupMode string `json:"setupMode,omitempty"`
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
