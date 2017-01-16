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
	ProductAPI      = "api"
	ProductVpnProxy = "vpnproxy"
	ProductGates    = "gates"
)

// AccountClient 包含针对账号 REST API 的各项操作
type AccountClient interface {
	// GetConfig 返回用于创建本 Client 实例的 AccountConfig
	GetConfig() (ret AccountConfig)

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

	// CreateAppGrant 将应用授权给用户
	CreateAppGrant(ctx context.Context, appURI, username string) (err error)

	// DeleteAppGrant 删除应用授权
	DeleteAppGrant(ctx context.Context, appURI, username string) (err error)

	// ListAppGrantedUsers 列出应用已授权的用户列表
	ListAppGrantedUsers(ctx context.Context, appURI string) (ret []AppGrantedUser, err error)

	// ListGrantedApps 列出已被授权的应用
	ListGrantedApps(ctx context.Context) (ret []AppInfo, err error)

	// GetGrantedAppKey 获取被授权应用的key
	GetGrantedAppKey(ctx context.Context, appURI string) (ret GrantedAppKey, err error)

	// ListGrants 获取自己授权给别人的应用列表
	ListGrants(ctx context.Context) (ret []GrantInfo, err error)

	// GetAppspecs 获得应用模板信息
	GetAppspecs(ctx context.Context, specURI string) (ret SpecInfo, err error)

	// ListPublicspecs 列出公开应用的模板
	ListPublicspecs(ctx context.Context) (ret []SpecInfo, err error)

	// ListGrantedspecs 列出被授权应用的模板
	ListGrantedspecs(ctx context.Context) (ret []SpecInfo, err error)

	// GetVendorManagedAppStatus 获得VendorManaged应用运行状态
	GetVendorManagedAppStatus(ctx context.Context, appURI string) (ret VendorManagedAppStatus, err error)

	// GetVendorManagedAppEntry 获得VendorManaged应用入口地址
	GetVendorManagedAppEntry(ctx context.Context, appURI string) (ret VendorManagedAppEntry, err error)

	// VendorManagedAppRepair 尝试修复VendorManaged应用
	VendorManagedAppRepair(ctx context.Context, appURI string) (err error)
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
	Title      string   `json:"title"`
	Region     string   `json:"region"`
	SpecURI    string   `json:"specUri"`
	SpecVer    uint32   `json:"specVer"`
	Privileges []string `json:"privileges"`
}

// AccountInfo 包含 Account 的相关信息
type AccountInfo struct {
	ID               uint32    `json:"id"`
	Name             string    `json:"name"`
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
	Privileges       []string  `json:"privileges,omitempty"`
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

// AppGrantedUser 包含列出应用被授权的用户信息
type AppGrantedUser struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

// GrantedAppKey 包含被授权应用的key信息
type GrantedAppKey struct {
	Ak string `json:"ak"`
	Sk string `json:"sk"`
}

// GrantInfo 应用授权信息
type GrantInfo struct {
	Account   string    `json:"account"`
	AppURI    string    `json:"appuri"`
	CreatedAt time.Time `json:"ctime"`
}

// SpecInfo 包含 Spec 的相关信息
type SpecInfo struct {
	URI        string    `json:"uri"`
	Owner      string    `json:"owner"`
	Title      string    `json:"title"`
	Ver        uint32    `json:"ver"`
	Verstr     string    `json:"verstr"`
	Desc       string    `json:"desc,omitempty"`
	Brief      string    `json:"brief"`
	Icon       string    `json:"icon"`
	Seedimg    string    `json:"seedimg"`
	Entryport  uint16    `json:"entryport"`
	Privileges []string  `json:"privileges"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
}

// VendorManagedAppStatus 包含应用运行状态信息
type VendorManagedAppStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// VendorManagedAppEntry 包含应用入口地址
type VendorManagedAppEntry struct {
	Entry string `json:"entry"`
}
