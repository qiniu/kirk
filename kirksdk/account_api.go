package kirksdk

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	// KeyStateEnabled 表示key的启用状态
	KeyStateEnabled = "enabled"

	// KeyStateDisabled 表示key的禁用状态
	KeyStateDisabled = "disabled"
)

const (
	// DefaultAccountHost 表示默认的账号 REST API 端点
	DefaultAccountHost = "https://app-api.qiniu.com"

	// DefaultIndexHost 表示默认的镜像 REST API 端点
	DefaultIndexHost = "https://index.qiniu.com"
)

const (
	ProductQcos         = "qcos"
	ProductQcosVpnProxy = "qcos_vpnproxy"
	ProductQcosGates    = "qcos_gates"
)

// AccountClient 包含针对账号 REST API 的各项操作
type AccountClient interface {

	// GetAccountInfo 用于得到 Account 的相关信息
	GetAccountInfo(ctx context.Context) (ret AccountInfo, err error)

	// CreateApp 用于在 Account 下创建 App
	CreateApp(ctx context.Context, appName string, args CreateAppArgs) (ret AppInfo, err error)

	// DeleteApp 用于删除 Account 下的 App
	DeleteApp(ctx context.Context, appURI string) (err error)

	// GetApp 用于得到 App 的相关信息
	GetApp(ctx context.Context, appName string) (ret AppInfo, err error)

	// GetAppKeys 用于得到 App 的 Key
	GetAppKeys(ctx context.Context, appURI string) (ret []KeyPair, err error)

	// ListApps 用于列出 Account 下的所有 App
	ListApps(ctx context.Context) (ret []AppInfo, err error)

	// ListManagedApps 用于列出 Account 下的所有 VendorManaged App
	ListManagedApps(ctx context.Context) (ret []AppInfo, err error)

	// GetRegion 用于得到 Region 相关的信息
	GetRegion(ctx context.Context, regionName string) (ret RegionInfo, err error)

	// ListRegions 用于列出所有可用的 Region
	ListRegions(ctx context.Context) (ret []RegionInfo, err error)

	// CreateAlertMethod 用于创建告警联系人
	CreateAlertMethod(ctx context.Context, appURI string, args CreateAlertMethodArgs) (ret AlertMethodInfo, err error)

	// DeleteAlertMethod 用于删除告警联系人
	DeleteAlertMethod(ctx context.Context, appURI string, id string) (err error)

	// GetAlertMethod 用于得到告警联系人相关信息
	GetAlertMethod(ctx context.Context, appURI string, id string) (ret AlertMethodInfo, err error)

	// ListAlertMethod 用于列出所有告警联系人
	ListAlertMethod(ctx context.Context, appURI string) (ret []AlertMethodInfo, err error)

	// UpdateAlertMethod 用于更新告警联系人
	UpdateAlertMethod(ctx context.Context, appURI string, id string, args UpdateAlertMethodArgs) (ret AlertMethodInfo, err error)

	// GetIndexClient 用于得到与镜像 REST API 交互的 IndexClient
	GetIndexClient(ctx context.Context) (client IndexClient, err error)

	// GetQcosClient 用于得到与某个 App 交互的 QcosClient
	GetQcosClient(ctx context.Context, appURI string) (client QcosClient, err error)
}

// AccountConfig 包含创建 AccountClient 所需的信息
type AccountConfig struct {
	AccessKey string `json:"ak"`
	SecretKey string `json:"sk"`
	Host      string `json:"appd_host"`
	UserAgent string
	Logger    *logrus.Logger
	Transport http.RoundTripper
}

// CreateAppArgs 包含创建一个 App 所需的信息
type CreateAppArgs struct {
	Title   string `json:"title"`
	Region  string `json:"region"`
	SpecURI string `json:"specUri"`
	SpecVer uint32 `json:"specVer"`
}

// AccountInfo 包含 Account 的相关信息
type AccountInfo struct {
	ID               uint32    `json:"id"`
	Name             string    `json:"uri"`
	Title            string    `json:"title"`
	CreationTime     time.Time `json:"ctime"`
	ModificationTime time.Time `json:"mtime"`
}

// AppInfo 包含 App 的相关信息
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

// AppExtendedInfo 包含 App 的相关扩展信息
type AppExtendedInfo struct {
	VendorURI string `json:"vendorUri,omitempty"`
	SpecURI   string `json:"specUri,omitempty"`
	SpecVer   uint32 `json:"specVer,omitempty"`
	SetupMode string `json:"setupMode,omitempty"`
}

// KeyPair 代表一对 Access Key， Secret Key 以及它们的启用状态
type KeyPair struct {
	AccessKey string `json:"ak"`
	SecretKey string `json:"sk"`
	State     string `json:"state"`
}

// RegionInfo 包含 Region 的相关信息
type RegionInfo struct {
	Name     string            `json:"name"`
	Desc     string            `json:"desc"`
	Products map[string]string `json:"products"`
}

// AlertMethods 代表若干个告警联系人
type AlertMethods struct {
	Methods []AlertMethodInfo `json:"methods"`
}

// AlertMethodInfo 包含告警联系人的相关信息
type AlertMethodInfo struct {
	ID          uint64 `json:"id"`
	Owner       uint32 `json:"owner"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	Nationality string `json:"nationality"`
	Code        string `json:"code"`
}

// CreateAlertMethodArgs 包含创建告警联系人所需的信息
type CreateAlertMethodArgs struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	Nationality string `json:"nationality"`
	Code        string `json:"code"`
}

// UpdateAlertMethodArgs 包含更新告警联系人所需的信息
type UpdateAlertMethodArgs struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	Nationality string `json:"nationality"`
	Code        string `json:"code"`
}
