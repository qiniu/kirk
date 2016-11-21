package kirksdk

import (
	"errors"
	"io"
	"time"

	"golang.org/x/net/context"
)

type QcosClient interface {
	// GET /v3/stacks
	ListStacks(ctx context.Context) (ret []StackInfo, err error)

	// POST /v3/stacks
	// Async
	CreateStack(ctx context.Context, args CreateStackArgs) (err error)
	// Sync
	SyncCreateStack(ctx context.Context, args CreateStackArgs) (err error)

	// POST /v3/stacks/<stackName>
	// Async
	UpdateStack(
		ctx context.Context, stackName string, args UpdateStackArgs) (err error)
	// Sync
	SyncUpdateStack(
		ctx context.Context, stackName string, args UpdateStackArgs) (err error)

	// GET /v3/stacks/<stackName>
	GetStack(ctx context.Context, stackName string) (ret StackInfo, err error)

	// GET /v3/stacks/<stackName>/export
	GetStackExport(
		ctx context.Context, stackName string) (ret CreateStackArgs, err error)

	// DELETE /v3/stacks/<stackName>
	// Async
	DeleteStack(ctx context.Context, stackName string) (err error)

	// POST /v3/stacks/<stackName>/start
	// Sync
	StartStack(ctx context.Context, stackName string) (err error)

	// POST /v3/stacks/<stackName>/stop
	// Sync
	StopStack(ctx context.Context, stackName string) (err error)

	// GET /v3/stacks/<stackName>/services
	ListServices(ctx context.Context, stackName string) (ret []ServiceInfo, err error)

	// POST /v3/stacks/<stackName>/services
	// Async
	CreateService(
		ctx context.Context, stackName string, args CreateServiceArgs) (err error)
	// Sync
	SyncCreateService(
		ctx context.Context, stackName string, args CreateServiceArgs) (err error)

	// GET /v3/stacks/<stackName>/services/<serviceName>/inspect
	GetServiceInspect(ctx context.Context,
		stackName string, serviceName string) (ret ServiceInfo, err error)

	// GET /v3/stacks/<stackName>/services/<serviceName>/export
	GetServiceExport(ctx context.Context,
		stackName string, serviceName string) (ret ServiceExportInfo, err error)

	// POST /v3/stacks/<stackName>/services/<serviceName>
	// Async
	UpdateService(ctx context.Context,
		stackName string, serviceName string, args UpdateServiceArgs) (err error)
	// Sync
	SyncUpdateService(ctx context.Context,
		stackName string, serviceName string, args UpdateServiceArgs) (err error)

	// POST /v3/stack/<stackName>/services/<serviceName>/deploy
	// Async
	DeployService(ctx context.Context,
		stackName string, serviceName string, args DeployServiceArgs) (err error)
	// Sync
	SyncDeployService(ctx context.Context,
		stackName string, serviceName string, args DeployServiceArgs) (err error)

	// POST /v3/stacks/<stackName>/services/<serviceName>/scale
	// Async
	ScaleService(ctx context.Context,
		stackName string, serviceName string, args ScaleServiceArgs) (err error)
	// Sync
	SyncScaleService(ctx context.Context,
		stackName string, serviceName string, args ScaleServiceArgs) (err error)

	// POST /v3/stacks/<stackName>/services/<serviceName>/start
	// Sync
	StartService(
		ctx context.Context, stackName string, serviceName string) (err error)

	// POST /v3/stacks/<stackName>/services/<serviceName>/stop
	// Sync
	StopService(
		ctx context.Context, stackName string, serviceName string) (err error)

	// DELETE /v3/stacks/<stackName>/services/<serviceName>
	// Async
	DeleteService(
		ctx context.Context, stackName string, serviceName string) (err error)

	// POST /v3/stacks/<stackName>/services/<serviceName>/volumes
	// Async
	CreateServiceVolume(ctx context.Context, stackName string,
		serviceName string, args CreateServiceVolumeArgs) (err error)
	// Sync
	SyncCreateServiceVolume(ctx context.Context, stackName string,
		serviceName string, args CreateServiceVolumeArgs) (err error)

	// POST /v3/stacks/<stackName>/services/<serviceName>/volumes/<volumeName>/extend
	// Async
	ExtendServiceVolume(ctx context.Context, stackName string,
		serviceName string, volumeName string, args ExtendVolumeArgs) (err error)
	// Sync
	SyncExtendServiceVolume(ctx context.Context, stackName string,
		serviceName string, volumeName string, args ExtendVolumeArgs) (err error)

	// DELETE /v3/stacks/<stackName>/services/<serviceName>/volumes/<volumeName>
	// Async
	DeleteServiceVolume(ctx context.Context, stackName string, serviceName string, volumeName string) (err error)
	// Sync
	SyncDeleteServiceVolume(ctx context.Context, stackName string, serviceName string, volumeName string) (err error)

	// POST /v3/stacks/<stackName>/services/<serviceName>/natip
	SetServiceNatIP(
		ctx context.Context, stackName string, serviceName string, args SetServiceNatIPArgs) (err error)

	// GET /v3/stacks/<stackName>/services/<serviceName>/natip
	GetServiceNatIP(
		ctx context.Context, stackName string, serviceName string) (natIP string, err error)

	// GET /v3/containers?stack=<stackName>&service=<serviceName>
	ListContainers(
		ctx context.Context, args ListContainersArgs) (ret []string, err error)

	// GET /v3/containers/<ip>/inspect
	GetContainerInspect(
		ctx context.Context, ip string) (ret ContainerInfo, err error)

	// POST /v3/containers/<ip>/start
	StartContainer(ctx context.Context, ip string) (err error)

	// POST /v3/containers/<ip>/stop
	// Sync
	StopContainer(ctx context.Context, ip string) (err error)

	// POST /v3/containers/<ip>/restart
	// Sync
	RestartContainer(ctx context.Context, ip string) (err error)

	// POST /v3/containers/<ip>/commit
	CommitContainerImage(
		ctx context.Context, ip string, args CommitContainerImageArgs) (err error)

	// POST /v3/containers/<ip>/exec
	ExecContainer(
		ctx context.Context, ip string, args ExecContainerArgs) (ret ExecContainerRet, err error)

	// POST /v3/containers/<ip>/exec/<execId>/resize
	ResizeContainerExecTerm(ctx context.Context,
		ip string, execID string, args ResizeContainerExecTermArgs) (err error)

	// POST /v3/containers/<ip>/exec/<execId>/start
	StartContainerExec(ctx context.Context,
		ip string, execID string, args StartContainerExecArgs, opts StartContainerExecOpts) (err error)

	// PUT /v3/containers/<ip>/webdav/files/<filePath>
	UploadToContainer(ctx context.Context,
		ip string, filePath string, rd io.Reader) (err error)

	// GET /v3/containers/<ip>/webdav/files/<filePath>
	DownloadFromContainer(ctx context.Context,
		ip string, filePath string) (rc io.ReadCloser, err error)

	// PROPFIND /v3/containers/<ip>/webdav/files/<filePath>
	StatContainerFile(ctx context.Context, ip string,
		filePath string, args StatContainerFileArgs) (rc io.ReadCloser, err error)

	// MKCOL /v3/containers/<ip>/webdav/files/<filePath>
	MkdirInContainer(ctx context.Context,
		ip string, filePath string) (err error)

	// GET /v3/logs/containers/<ip>/realtime?since=<since>&tail=<tail>
	GetContainerLogsRealtime(ctx context.Context,
		ip, since, tail string, opts GetContainerLogsRealtimeOpts) (stream io.ReadCloser, err error)

	// GET /v3/logs/search/<repoType>?q=<query>&from=<from>&size=<size>&sort=<sort>
	SearchContainerLogs(ctx context.Context, args SearchContainerLogsArgs) (res LogsSearchResult, err error)

	// AP(Access Point) APIs

	// GET /v3/aps | /v3/aps?stack=<stack> | GET /v3/aps?service=<service>
	ListAps(ctx context.Context, args ListApsArgs) (ret []ListApInfo, err error)

	// POST /v3/aps
	CreateAp(ctx context.Context, args CreateApArgs) (ret ListApInfo, err error)

	// GET  /v3/aps/search?ip=<IP> | GET  /v3/aps/search?domain=<domain>
	SearchAp(ctx context.Context, mode string, searchArg string) (ret FullApInfo, err error)

	// GET  /v3/aps/<apid>
	GetAp(ctx context.Context, apid string) (ret FullApInfo, err error)

	// POST /v3/aps/<apid>
	UpdateAp(ctx context.Context, apid string, args SetApDescArgs) (err error)

	// POST /v3/aps/<apid>/<port>
	SetApPort(ctx context.Context, apid string, port string, args SetApPortArgs) (err error)

	// DELETE /v3/aps/<apid>/<port>
	DeleteApPort(ctx context.Context, apid string, port string) (err error)

	// POST /v3/aps/<apid>/portrange/<from>/<to>
	SetApPortRange(
		ctx context.Context, apid string, fromPort string, toPort string, args SetApPortRangeArgs) (err error)

	//DELETE /v3/aps/<apid>/<portrange>/<from>/<to>
	DeleteApPortRange(
		ctx context.Context, apid string, fromPort string, toPort string) (err error)

	// DELETE /v3/aps/<apid>
	DeleteAp(ctx context.Context, apid string) (err error)

	// GET  /v3/aps/<apid>/<port>/healthcheck
	GetHealthcheck(ctx context.Context, apid string, port string) (ret map[string]string, err error)

	// POST /v3/aps/<apid>/<port>/setcontainer
	ApSetContainer(ctx context.Context, apid string, port string, args []SetApContainerOptionsArgs) (err error)

	// POST /v3/aps/<apid>/publish
	PublishUserDomain(ctx context.Context, apid string, args SetUserDomainArgs) (err error)

	// POST /v3/aps/<apid>/unpublish
	UnpublishUserDomain(ctx context.Context, apid string, args SetUserDomainArgs) (err error)

	// GET /v3/aps/providers
	ListProviders(ctx context.Context) (ret []string, err error)

	// GET /v3/jobs
	ListJobs(ctx context.Context) (ret []JobInfo, err error)

	// GET /v3/jobs/<name>
	GetJob(ctx context.Context, name string) (ret JobInfo, err error)

	// DELETE /v3/jobs/<name>
	DeleteJob(ctx context.Context, name string) (err error)

	// POST /v3/jobs
	CreateJob(ctx context.Context, args CreateJobArgs) (err error)

	// POST /v3/jobs/<name>
	UpdateJob(ctx context.Context, name string, args UpdateJobArgs) (err error)

	// POST /v3/jobs/<name>/run
	RunJob(ctx context.Context, name string, args RunJobArgs) (ret JobInstanceID, err error)

	// GET /v3/jobs/<name>/instances/<id>
	GetJobInstance(ctx context.Context, name string, id string) (ret JobInstance, err error)

	// DELETE /v3/jobs/<name>/instances/<id>
	DeleteJobInstance(ctx context.Context, name string, id string) (err error)

	// POST /v3/jobs/<name>/instances/<id>/stop
	StopJobInstance(ctx context.Context, name string, id string) (err error)

	//  POST /v3/alert/aps/<apid>
	UpdateApAlert(ctx context.Context, apid string, args UpdateApAlertArgs) (err error)

	// DELETE /v3/alert/aps/<apid>
	DeleteApAlert(ctx context.Context, apid string, level string) (err error)

	//  GET /v3/alert/aps/<apid>?level=<level>
	GetApAlert(ctx context.Context, apid string, level string) (ret []ApAlertInfo, err error)

	// POST /v3/alert/stacks/<stackName>/services/<serviceName>
	UpdateServiceAlert(ctx context.Context, stack, service string, args UpdateContainerAlertArgs) (err error)

	// POST /v3/alert/stacks/<stackName>/services/<serviceName>/all
	UpdateAllContainerAlert(ctx context.Context, stack, service string, args UpdateContainerAlertArgs) (err error)

	// POST /v3/alert/containers/<ip>
	UpdateContainerAlert(ctx context.Context, ip string, args UpdateContainerAlertArgs) (err error)

	// DELETE /v3/alert/stacks/<stackName>/services/<serviceName>
	DeleteServiceAlert(ctx context.Context, stack, service string, level string) (err error)

	// DELETE /v3/alert/containers/<ip>
	DeleteContainerAlert(ctx context.Context, ip string, level string) (err error)

	// GET /v3/alert/stacks/<stackName>/services/<serviceName>?level=<level>
	GetServiceAlert(ctx context.Context, stack, service string, level string) (ret []ContainerAlertInfo, err error)

	// GET /v3/alert/containers/<ip>?level=<level>
	GetContainerAlert(ctx context.Context, ip string, level string) (ret []ContainerAlertInfo, err error)
}

const (
	StatusRunning       = Status("RUNNING")
	StatusPartlyRunning = Status("PARTIALLY-RUNNING")
	StatusNotRunning    = Status("NOT-RUNNING")
	StatusFault         = Status("FAULT")
)

const (
	StateCreate         = State("CREATING")
	StateScaling        = State("SCALING")
	StateAutoUpdating   = State("AUTO-UPDATING")
	StateManualUpdating = State("MANUAL-UPDATING")
	StateStarting       = State("STARTING")
	StateStopping       = State("STOPPING")
	StateStopped        = State("STOPPED")
	StateDeployed       = State("DEPLOYED")
)

const (
	ApBackendDefaultWeight = 10000

	ApTypePublicIPStr  = "PUBLIC_IP"
	ApTypePrivateIPStr = "INTERNAL_IP"
	ApTypeDomainStr    = "DOMAIN"
)

type Status string

type State string

type StackInfo struct {
	IsDeployed bool     `json:"isDeployed"`
	Metadata   []string `json:"metadata"`
	Name       string   `json:"name"`
	Services   []string `json:"services"`
	Status     Status   `json:"status"`
}

type CreateStackArgs struct {
	Metadata []string            `json:"metadata"`
	Name     string              `json:"name"`
	Services []CreateServiceArgs `json:"services"`
}

type UpdateStackArgs struct {
	Metadata []string `json:"metadata"`
}

type CreateServiceArgs struct {
	InstanceNum       int          `json:"instanceNum"`
	UpdateParallelism int          `json:"updateParallelism"`
	Metadata          []string     `json:"metadata"`
	Name              string       `json:"name"`
	Spec              ServiceSpec  `json:"spec"`
	Stateful          bool         `json:"stateful"`
	Volumes           []VolumeSpec `json:"volumes"`
}

type VolumeSpec struct {
	FsType    string `json:"fsType"`
	MountPath string `json:"mountPath"`
	Name      string `json:"name"`
	UnitType  string `json:"unitType"`
}

type ServiceInfo struct {
	ContainerIPs      []string          `json:"containerIps"`
	InstanceNum       int               `json:"instanceNum"`
	UpdateParallelism int               `json:"updateParallelism"`
	Metadata          []string          `json:"metadata"`
	Name              string            `json:"name"`
	Revision          int               `json:"revision"`
	Spec              ServiceSpecExport `json:"spec"`
	Stack             string            `json:"stack"`
	State             State             `json:"state"`
	Stateful          bool              `json:"stateful"`
	Status            Status            `json:"status"`
	UpdateSpec        ServiceSpecExport `json:"updateSpec"`
	Volumes           []VolumeSpec      `json:"volumes"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
}

type ServiceExportInfo struct {
	InstanceNum       int               `json:"instanceNum"`
	UpdateParallelism int               `json:"updateParallelism"`
	Metadata          []string          `json:"metadata"`
	Name              string            `json:"name"`
	Spec              ServiceSpecExport `json:"spec"`
	Stateful          bool              `json:"stateful"`
	Volumes           []VolumeSpec      `json:"volumes"`
}

type ServiceSpecExport struct {
	AutoRestart   string             `json:"autoRestart"`
	Command       []string           `json:"command"`
	EntryPoint    []string           `json:"entryPoint"`
	Envs          []string           `json:"envs"`
	Hosts         []string           `json:"hosts"`
	Image         string             `json:"image"`
	LogCollectors []LogCollectorSpec `json:"logCollectors"`
	StopGraceSec  int                `json:"stopGraceSec"`
	WorkDir       string             `json:"workDir"`
	UnitType      string             `json:"unitType"`
}

// If empty, service will use default values.
type ServiceSpec struct {
	AutoRestart   string             `json:"autoRestart,omitempty"`
	Command       []string           `json:"command,omitempty"`
	EntryPoint    []string           `json:"entryPoint,omitempty"`
	Envs          []string           `json:"envs,omitempty"`
	Hosts         []string           `json:"hosts,omitempty"`
	Image         string             `json:"image,omitempty"`
	LogCollectors []LogCollectorSpec `json:"logCollectors,omitempty"`
	StopGraceSec  int                `json:"stopGraceSec,omitempty"`
	WorkDir       string             `json:"workDir,omitempty"`
	UnitType      string             `json:"unitType,omitempty"`
	GpuUUIDs      []string           `json:"gpuUUIDs,omitempty"` // do not use. prepare for gpu service
}

type LogCollectorSpec struct {
	Directory string   `json:"directory"`
	Patterns  []string `json:"patterns"`
}

type UpdateServiceArgs struct {
	ManualUpdate      bool        `json:"manualUpdate"`
	Metadata          []string    `json:"metadata"`
	Spec              ServiceSpec `json:"spec"`
	UpdateParallelism int         `json:"updateParallelism"`
}

type DeployServiceArgs struct {
	Operation string `json:"operation"`
}

type ScaleServiceArgs struct {
	InstanceNum int `json:"instanceNum"`
}

type ExtendVolumeArgs struct {
	UnitType string `json:"unitType"`
}

type CreateServiceVolumeArgs struct {
	VolumeSpec
}

type SetServiceNatIPArgs struct {
	NatIP string `json:"natip"`
}

type ListContainersArgs struct {
	StackName   string `json:"stackName"`
	ServiceName string `json:"serviceName"`
}

type ContainerInfo struct {
	CPU struct {
		CoreUsage  float64 `json:"coreUsage"`
		TotalUsage float64 `json:"totalUsage"`
	} `json:"cpu"`
	Disk struct {
		Usage int `json:"usage"`
	} `json:"disk"`
	ExitCode int    `json:"exitCode"`
	ExitMsg  string `json:"exitMsg"`
	IP       string `json:"ip"`
	Memory   struct {
		Cache int `json:"cache"`
		Usage int `json:"usage"`
	} `json:"memory"`
	Network struct {
		RxBs    float64 `json:"rxBs"`
		RxBytes int     `json:"rxBytes"`
		TxBs    float64 `json:"txBs"`
		TxBytes int     `json:"txBytes"`
	} `json:"network"`
	Revision   int                        `json:"revision"`
	Service    string                     `json:"service"`
	Stack      string                     `json:"stack"`
	Status     Status                     `json:"status"`
	Volumes    map[string]VolumeUsageInfo `json:"volumes,omitempty"`
	GpuUUIDs   []string                   `json:"gpuUUIDs,omitempty"`
	CreatedAt  time.Time                  `json:"createdAt"`
	StartedAt  time.Time                  `json:"startedAt"`
	FinishedAt time.Time                  `json:"finishedAt"`
}

type VolumeUsageInfo struct {
	Iops int `json:"iops"`
	Size struct {
		Total int `json:"total"`
		Used  int `json:"used"`
	}
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CommitContainerImageArgs struct {
	Image string `json:"image"`
}

type ExecContainerArgs struct {
	Command []string `json:"command"`
}

type ResizeContainerExecTermArgs struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type StartContainerExecArgs struct {
	Mode string `json:"mode"`
}

type CreateApArgs struct {
	Type      string `json:"type"`
	Provider  string `json:"provider"`
	Bandwidth int    `json:"bandwidthMbps"`
	UnitType  string `json:"unitType"`
	Host      string `json:"host"`
	Title     string `json:"title"`
}

type ListApsArgs struct {
	Service string `json:"service"`
	Stack   string `json:"stack"`
	Title   string `json:"title"`
}

type ListApInfo struct {
	ApID      string   `json:"apid"`
	Type      string   `json:"type"`
	Title     string   `json:"title"`
	Bandwidth int      `json:"bandwidthMbps"`
	IP        string   `json:"ip,omitempty"`
	Domains   []string `json:"domains,omitempty"`
	Provider  string   `json:"provider"`
	Host      string   `json:"host,omitempty"`
	UnitType  string   `json:"unitType,omitempty"`
}

type ApProxyOpts struct {
	FailTimeoutMS     int      `json:"failTimeoutMs"`
	MaxFails          int      `json:"maxFails"`
	NextUpstreamTries int      `json:"nextUpstreamTries"`
	NextUpstreamCond  []string `json:"nextUpstreamCond"`
}

var DefaultApProxyOpts = ApProxyOpts{
	FailTimeoutMS:     10000,
	MaxFails:          1,
	NextUpstreamTries: 3,
	NextUpstreamCond:  []string{"error"},
}

type ApHealthCheckOpts struct {
	Enabled     bool   `json:"enabled"`
	Path        string `json:"path"`
	HttpOkCodes []int  `json:"httpOkCodes"`
}

type ApPortInfo struct {
	Proto           string            `json:"proto"`
	FPort           string            `json:"frontendPort"`
	BPort           string            `json:"backendPort"`
	GroupID         int               `json:"groupId"`
	SessionTmoSec   int               `json:"sessionTimeoutSec"`
	ProxyOpts       ApProxyOpts       `json:"proxyOptions"`
	HealthCheckOpts ApHealthCheckOpts `json:"healthCheck"`
	Backends        []struct {
		Stack         string `json:"stack"`
		Service       string `json:"service"`
		DefaultWeight int    `json:"weight"`
		ActualWeight  int    `json:"actualWeight"`
	} `json:"backends"`
	ContainerOptions []struct {
		CotainerIP string  `json:"ip"`
		Ratio      float64 `json:"ratio"`
	} `json:"containerOptions"`
}

type FullApInfo struct {
	ApID        int          `json:"apid"`
	Type        string       `json:"type"`
	Title       string       `json:"title"`
	IP          string       `json:"ip,omitempty"`
	Domain      string       `json:"domain,omitempty"`
	Provider    string       `json:"provider"`
	Bandwidth   int          `json:"bandwidthMbps"`
	Traffic     int          `json:"trafficBytes"`
	UserDomains []string     `json:"userDomains,omitempty"`
	Host        string       `json:"host,omitempty"`
	UnitType    string       `json:"unitType,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	Ports       []ApPortInfo `json:"ports"`
}

type SetApDescArgs struct {
	UnitType  string `json:"unitType"`
	Host      string `json:"host"`
	Title     string `json:"title"`
	Bandwidth int    `json:"bandwidthMbps"`
}

type ApBackendArgs struct {
	Stack   string `json:"stack"`
	Service string `json:"service"`
	Weight  int    `json:"weight"`
}

type SetApPortArgs struct {
	Proto         string `json:"proto"`
	BackendPort   int    `json:"backendPort"`
	SessionTmoSec int    `json:"sessionTimeoutSec"`

	ProxyOpt    *ApProxyOpts       `json:"proxyOptions"`
	HealthCheck *ApHealthCheckOpts `json:"healthCheck"`
	Backends    []ApBackendArgs    `json:"backends"`
}

type SetApPortRangeArgs struct {
	Proto             string          `json:"proto"`
	SessionTimeoutSec int             `json:"sessionTimeoutSec"`
	Backends          []ApBackendArgs `json:"backends"`
}
type SetUserDomainArgs struct {
	UserDomain string `json:"userDomain"`
}

type SetApContainerOptionsArgs struct {
	IP    string  `json:"ip"`
	Ratio float64 `json:"ratio"`
}

type ExecContainerRet struct {
	ExecID string `json:"execId"`
}

type StartContainerExecOpts struct {
	InStream  io.Reader
	OutStream io.Writer
	ErrStream io.Writer
	ReadyCh   chan struct{}
	ErrorCh   chan error
}

type StatContainerFileArgs struct {
	Depth int
}

type GetContainerLogsRealtimeOpts struct {
	ExitCh  chan struct{}
	ErrorCh chan error
}

type SearchContainerLogsArgs struct {
	RepoType string
	Query    string
	From     int
	Size     int
	Sort     string
}

type LogsSearchResult struct {
	Total int   `json:"total"`
	Data  []Hit `json:"data"`
}

type Hit struct {
	Log         string    `json:"log"`
	CollectedAt time.Time `json:"collectedAt"`
	PodIP       string    `json:"podIp"`
	ProcessName string    `json:"processName"`
	GateID      string    `json:"gateId"`
	Domain      string    `json:"domain"`
}

var (
	ErrNotImplement  = errors.New("not support")
	ErrNoSuchProcess = errors.New("no such process")
	ErrNoSuchExec    = errors.New("no such exec")
	ErrResultError   = errors.New("result error")
	ErrNoSuchEntry   = errors.New("no such entry")
)

type JobTaskSpec struct {
	Image         string             `json:"image,omitempty"`
	Command       []string           `json:"command,omitempty"`
	EntryPoint    []string           `json:"entryPoint,omitempty"`
	Envs          []string           `json:"envs,omitempty"`
	Hosts         []string           `json:"hosts,omitempty"`
	LogCollectors []LogCollectorSpec `json:"logCollectors,omitempty"`
	WorkDir       string             `json:"workDir,omitempty"`
	InstanceNum   int                `json:"instanceNum,omitempty"`
	UnitType      string             `json:"unitType,omitempty"`
}

type JobInfo struct {
	Metadata []string               `json:"metadata"`
	Name     string                 `json:"name"`
	Spec     map[string]JobTaskSpec `json:"spec"`
	Revision int                    `json:"revision"`

	RunAt   string `json:"runAt"`
	Timeout int    `json:"timeout"`
	Mode    string `json:"mode"`

	LastRun    time.Time `json:"lastRun"`
	LastStatus JobStatus `json:"lastStatus"`

	JobInstances []string `json:"jobInstances"`

	Created time.Time `db:"created" json:"created"`
	Updated time.Time `db:"updated" json:"updated"`
}

type CreateJobArgs struct {
	Name     string                 `json:"name"`
	Spec     map[string]JobTaskSpec `json:"spec"`
	Mode     string                 `json:"mode"`
	Metadata []string               `json:"metadata,omitempty"`
	RunAt    string                 `json:"runAt,omitempty"`
	Timeout  int                    `json:"timeout,omitempty"`
}

type UpdateJobArgs struct {
	Spec     map[string]JobTaskSpec `json:"spec,omitempty"`
	Metadata []string               `json:"metadata,omitempty"`
	RunAt    string                 `json:"runAt,omitempty"`
	Timeout  int                    `json:"timeout,omitempty"`
	Mode     string                 `json:"mode,omitempty"`
}

type RunJobArgs struct {
	Spec map[string]JobTaskSpecEx `json:"spec,omitempty"`
}

type JobInstanceID struct {
	ID string `json:"jobInstanceId"`
}

type JobTaskSpecEx struct {
	WorkDir       string             `json:"workDir,omitempty"`
	LogCollectors []LogCollectorSpec `json:"logCollectors,omitempty"`
	Command       []string           `json:"command,omitempty"`
	EntryPoint    []string           `json:"entryPoint,omitempty"`
	Envs          []string           `json:"envs,omitempty"`
	Hosts         []string           `json:"hosts,omitempty"`
	UnitType      string             `json:"unitType,omitempty"`
	Deps          []string           `json:"deps,omitempty"`
	InstanceNum   int                `json:"instanceNum,omitempty"`
}

const (
	JobStatusRunning = JobStatus("RUNNING")
	JobStatusWaiting = JobStatus("WAITING")
	JobStatusStopped = JobStatus("STOPPED")
	JobStatusFail    = JobStatus("FAIL")
	JobStatusSuccess = JobStatus("SUCCESS")
	JobStatusTimeout = JobStatus("TIMEOUT")
)

type JobStatus string

type JobTask struct {
	ID         string    `json:"taskId"`
	Name       string    `json:"name"`
	Index      int       `json:"index"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`
}

type JobInstance struct {
	InstanceID string                 `json:"instanceId"`
	Name       string                 `json:"name"`
	Spec       map[string]JobTaskSpec `json:"spec"`
	Status     JobStatus              `json:"status"`
	Tasks      []JobTask              `json:"tasks"`
	CreatedAt  time.Time              `json:"createdAt"`
	StartedAt  time.Time              `json:"startedAt"`
	FinishedAt time.Time              `json:"finishedAt"`
}

type AlertLevelArgs struct {
	Level string `json:"level,omitempty"`
}

type AlertMethod struct {
	ID     uint64 `json:"id"`
	Email  string `json:"email,omitempty"`
	Mobile string `json:"mobile,omitempty"`
	Code   string `json:"code,omitempty"`
}

type AlertApThreshold struct {
	Uplink   string `json:"uplink,omitempty"`
	Downlink string `json:"downlink,omitempty"`
	Qps      string `json:"qps,omitempty"`
}

type UpdateApAlertArgs struct {
	Level     string            `json:"level"`
	Threshold *AlertApThreshold `json:"threshold"`
	Methods   []*AlertMethod    `json:"methods"`
}

type ApAlertInfo struct {
	Level     string            `json:"level"`
	Threshold *AlertApThreshold `json:"threshold"`
	Methods   []*AlertMethod    `json:"methods"`
}

type AlertContainerThreshold struct {
	CPU    float64 `json:"cpu,omitempty"`
	Mem    float64 `json:"mem,omitempty"`
	Rootfs float64 `json:"rootfs,omitempty"`
	Volume float64 `json:"volume,omitempty"`
	NetRx  float64 `json:"netrx,omitempty"`
	NetTx  float64 `json:"nettx,omitempty"`
}

type UpdateContainerAlertArgs struct {
	Level     string                   `json:"level"`
	Threshold *AlertContainerThreshold `json:"threshold"`
	Methods   []*AlertMethod           `json:"methods"`
}

type ContainerAlertInfo struct {
	Level     string                   `json:"level"`
	Threshold *AlertContainerThreshold `json:"threshold"`
	Methods   []*AlertMethod           `json:"methods"`
}
