package kirksdk

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"qiniupkg.com/kirk/kirksdk/mac"
	"qiniupkg.com/x/rpc.v7"
)

const DefaultStack = "default"
const waitTimeout = 120 * time.Second

const MultiStatus = 207

type QcosConfig struct {
	AccessKey string // OPTIONAL assume client inside qcos if not set
	SecretKey string // OPTIONAL assume client inside qcos if not set
	Host      string
	UserAgent string
	Transport http.RoundTripper
	Logger    *logrus.Logger
}

type qcosClientImp struct {
	host   string
	logger *logrus.Logger
	client rpc.Client
	kmac   *mac.Mac
}

func NewQcosClient(cfg QcosConfig) QcosClient {

	p := new(qcosClientImp)

	p.host = cleanHost(cfg.Host)

	if cfg.Logger == nil {
		cfg.Logger = logrus.New()
	}
	p.logger = cfg.Logger

	cfg.Transport = newKirksdkTransport(cfg.UserAgent, cfg.Transport)

	if cfg.AccessKey == "" { // client used inside intranet
		p.client = rpc.Client{&http.Client{Transport: cfg.Transport}}
	} else {
		p.kmac = mac.New(cfg.AccessKey, cfg.SecretKey)
		p.client = rpc.Client{mac.NewClient(p.kmac, cfg.Transport)}
	}

	return p
}

// GET /v3/stacks
func (p *qcosClientImp) ListStacks(ctx context.Context) (ret []StackInfo, err error) {

	url := fmt.Sprintf("%s/v3/stacks", p.host)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// POST /v3/stacks
func (p *qcosClientImp) CreateStack(
	ctx context.Context, args CreateStackArgs) (err error) {

	url := fmt.Sprintf("%s/v3/stacks", p.host)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

func (p *qcosClientImp) SyncCreateStack(
	ctx context.Context, args CreateStackArgs) (err error) {
	err = p.CreateStack(ctx, args)
	if err != nil {
		return
	}
	err = p.wait4StackRunning(args.Name, waitTimeout)
	if err != nil {
		return
	}
	return
}

// POST /v3/stacks/<stackName>
func (p *qcosClientImp) UpdateStack(ctx context.Context, stackName string,
	args UpdateStackArgs) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s", p.host, stackName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

func (p *qcosClientImp) SyncUpdateStack(ctx context.Context, stackName string,
	args UpdateStackArgs) (err error) {
	err = p.UpdateStack(ctx, stackName, args)
	if err != nil {
		return
	}
	err = p.wait4StackRunning(stackName, waitTimeout)
	if err != nil {
		return
	}
	return
}

// GET /v3/stacks/<stackName>
func (p *qcosClientImp) GetStack(
	ctx context.Context, stackName string) (ret StackInfo, err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s", p.host, stackName)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// GET /v3/stacks/<stackName>/export
func (p *qcosClientImp) GetStackExport(
	ctx context.Context, stackName string) (ret CreateStackArgs, err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/export", p.host, stackName)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// DELETE /v3/stacks/<stackName>
func (p *qcosClientImp) DeleteStack(
	ctx context.Context, stackName string) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s", p.host, stackName)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

// POST /v3/stacks/<stackName>/start
func (p *qcosClientImp) StartStack(ctx context.Context, stackName string) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/start", p.host, stackName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, nil)
	return
}

// POST /v3/stacks/<stackName>/stop
func (p *qcosClientImp) StopStack(ctx context.Context, stackName string) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/stop", p.host, stackName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, nil)
	return
}

// GET /v3/stacks/<stackName>/services
func (p *qcosClientImp) ListServices(
	ctx context.Context, stackName string) (ret []ServiceInfo, err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services", p.host, stackName)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// POST /v3/stacks/<stackName>/services
func (p *qcosClientImp) CreateService(
	ctx context.Context, stackName string, args CreateServiceArgs) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services", p.host, stackName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

func (p *qcosClientImp) SyncCreateService(
	ctx context.Context, stackName string, args CreateServiceArgs) (err error) {
	err = p.CreateService(ctx, stackName, args)
	if err != nil {
		return
	}
	err = p.wait4ServiceRunning(stackName, args.Name, waitTimeout)
	if err != nil {
		return
	}
	return
}

// GET /v3/stacks/<stackName>/services/<serviceName>/inspect
func (p *qcosClientImp) GetServiceInspect(ctx context.Context,
	stackName string, serviceName string) (ret ServiceInfo, err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s/inspect", p.host, stackName, serviceName)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// GET /v3/stacks/<stackName>/services/<serviceName>/export
func (p *qcosClientImp) GetServiceExport(ctx context.Context, stackName string,
	serviceName string) (ret ServiceExportInfo, err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s/export", p.host, stackName, serviceName)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// POST /v3/stacks/<stackName>/services/<serviceName>
func (p *qcosClientImp) UpdateService(ctx context.Context, stackName string,
	serviceName string, args UpdateServiceArgs) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s", p.host, stackName, serviceName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

func (p *qcosClientImp) SyncUpdateService(ctx context.Context, stackName string,
	serviceName string, args UpdateServiceArgs) (err error) {
	err = p.UpdateService(ctx, stackName, serviceName, args)
	if err != nil {
		return
	}

	if args.ManualUpdate == false {
		err = p.wait4ServiceRunning(stackName, serviceName, waitTimeout)
		if err != nil {
			return
		}
	}
	return
}

// POST /v3/stacks/<stackName>/services/<serviceName>/deploy
func (p *qcosClientImp) DeployService(ctx context.Context,
	stackName string, serviceName string, args DeployServiceArgs) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s/deploy", p.host, stackName, serviceName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

func (p *qcosClientImp) SyncDeployService(ctx context.Context,
	stackName string, serviceName string, args DeployServiceArgs) (err error) {
	err = p.DeployService(ctx, stackName, serviceName, args)
	if err != nil {
		return
	}
	if args.Operation == "ROLLBACK" || args.Operation == "COMPLETE" {
		err = p.wait4ServiceRunning(stackName, serviceName, waitTimeout)
		if err != nil {
			return
		}
	}
	op := strings.Split(args.Operation, " ")
	switch op[0] {
	case "COMPLETE", "ROLLBACK":
		err = p.wait4ServiceRunning(stackName, serviceName, waitTimeout)
		if err != nil {
			return
		}
	default:
		return
	}
	return
}

// POST /v3/stacks/<stackName>/services/<serviceName>/scale
func (p *qcosClientImp) ScaleService(ctx context.Context,
	stackName string, serviceName string, args ScaleServiceArgs) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s/scale", p.host, stackName, serviceName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

func (p *qcosClientImp) SyncScaleService(ctx context.Context,
	stackName string, serviceName string, args ScaleServiceArgs) (err error) {
	err = p.ScaleService(ctx, stackName, serviceName, args)
	if err != nil {
		return
	}
	err = p.wait4ServiceRunning(stackName, serviceName, waitTimeout)
	if err != nil {
		return
	}
	return
}

// POST /v3/stacks/<stackName>/services/<serviceName>/start
func (p *qcosClientImp) StartService(
	ctx context.Context, stackName string, serviceName string) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s/start", p.host, stackName, serviceName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, nil)
	return
}

func (p *qcosClientImp) SyncStartService(
	ctx context.Context, stackName string, serviceName string) (err error) {
	err = p.StartService(ctx, stackName, serviceName)
	if err != nil {
		return
	}
	err = p.wait4ServiceRunning(stackName, serviceName, waitTimeout)
	if err != nil {
		return
	}
	return
}

// POST /v3/stacks/<stackName>/services/<serviceName>/stop
func (p *qcosClientImp) StopService(
	ctx context.Context, stackName string, serviceName string) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s/stop", p.host, stackName, serviceName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, nil)
	return
}

func (p *qcosClientImp) SyncStopService(
	ctx context.Context, stackName string, serviceName string) (err error) {
	err = p.StopService(ctx, stackName, serviceName)
	if err != nil {
		return
	}
	err = p.wait4ServiceStopped(stackName, serviceName, waitTimeout)
	if err != nil {
		return
	}
	return
}

// DELETE /v3/stacks/<stackName>/services/<serviceName>
func (p *qcosClientImp) DeleteService(ctx context.Context, stackName string, serviceName string) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s", p.host, stackName, serviceName)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

// POST /v3/stacks/<stackName>/services/<serviceName>/volumes
func (p *qcosClientImp) CreateServiceVolume(ctx context.Context, stackName string,
	serviceName string, args CreateServiceVolumeArgs) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s/volumes", p.host, stackName, serviceName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

func (p *qcosClientImp) SyncCreateServiceVolume(ctx context.Context, stackName string,
	serviceName string, args CreateServiceVolumeArgs) (err error) {
	err = p.CreateServiceVolume(ctx, stackName, serviceName, args)
	if err != nil {
		return
	}
	err = p.wait4ServiceRunning(stackName, serviceName, waitTimeout)
	if err != nil {
		return
	}
	return
}

// POST /v3/stacks/<stackName>/services/<serviceName>/volumes/<volumeName>/extend
func (p *qcosClientImp) ExtendServiceVolume(ctx context.Context, stackName string,
	serviceName string, volumeName string, args ExtendVolumeArgs) (err error) {

	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s/volumes/%s/extend", p.host, stackName, serviceName, volumeName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

func (p *qcosClientImp) SyncExtendServiceVolume(ctx context.Context, stackName string,
	serviceName string, volumeName string, args ExtendVolumeArgs) (err error) {
	err = p.ExtendServiceVolume(ctx, stackName, serviceName, volumeName, args)
	if err != nil {
		return
	}
	err = p.wait4ServiceRunning(stackName, serviceName, waitTimeout)
	if err != nil {
		return
	}
	return
}

// DELETE /v3/stacks/<stackName>/services/<serviceName>/volumes/<volumeName>
func (p *qcosClientImp) DeleteServiceVolume(ctx context.Context, stackName string, serviceName string, volumeName string) (err error) {
	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s/volumes/%s", p.host, stackName, serviceName, volumeName)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

func (p *qcosClientImp) SyncDeleteServiceVolume(ctx context.Context, stackName string, serviceName string, volumeName string) (err error) {
	err = p.DeleteServiceVolume(ctx, stackName, serviceName, volumeName)
	if err != nil {
		return
	}
	err = p.wait4ServiceRunning(stackName, serviceName, waitTimeout)
	if err != nil {
		return
	}
	return
}

//POST /v3/stacks/<stackName>/services/<serviceName>/natip
func (p *qcosClientImp) SetServiceNatIP(ctx context.Context, stackName string, serviceName string, args SetServiceNatIPArgs) (err error) {
	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s/natip", p.host, stackName, serviceName)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

//GET /v3/stacks/<stackName>/services/<serviceName>/natip
func (p *qcosClientImp) GetServiceNatIP(ctx context.Context, stackName string, serviceName string) (natIP string, err error) {
	if stackName == "" {
		stackName = DefaultStack
	}

	url := fmt.Sprintf("%s/v3/stacks/%s/services/%s/natip", p.host, stackName, serviceName)
	err = p.client.Call(ctx, &natIP, "GET", url)
	return
}

// GET /v3/containers?stack=<stackName>&service=<serviceName>
func (p *qcosClientImp) ListContainers(
	ctx context.Context, args ListContainersArgs) (ret []string, err error) {

	queryString := ""
	if args.StackName != "" {
		queryString = fmt.Sprintf("?stack=%s", args.StackName)
	}
	if args.ServiceName != "" {
		if queryString == "" {
			queryString = fmt.Sprintf("?service=%s", args.ServiceName)
		} else {
			queryString += fmt.Sprintf("&service=%s", args.ServiceName)
		}
	}
	url := fmt.Sprintf("%s/v3/containers%s", p.host, queryString)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// GET /v3/containers/<ip>/inspect
func (p *qcosClientImp) GetContainerInspect(
	ctx context.Context, ip string) (ret ContainerInfo, err error) {

	url := fmt.Sprintf("%s/v3/containers/%s/inspect", p.host, ip)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// POST /v3/containers/<ip>/start
func (p *qcosClientImp) StartContainer(ctx context.Context, ip string) (err error) {

	url := fmt.Sprintf("%s/v3/containers/%s/start", p.host, ip)
	err = p.client.CallWithJson(ctx, nil, "POST", url, nil)
	return
}

// POST /v3/containers/<ip>/stop
func (p *qcosClientImp) StopContainer(ctx context.Context, ip string) (err error) {

	url := fmt.Sprintf("%s/v3/containers/%s/stop", p.host, ip)
	err = p.client.CallWithJson(ctx, nil, "POST", url, nil)
	return
}

// POST /v3/containers/<ip>/restart
func (p *qcosClientImp) RestartContainer(ctx context.Context, ip string) (err error) {

	url := fmt.Sprintf("%s/v3/containers/%s/restart", p.host, ip)
	err = p.client.CallWithJson(ctx, nil, "POST", url, nil)
	return
}

// POST /v3/containers/<ip>/commit
func (p *qcosClientImp) CommitContainerImage(
	ctx context.Context, ip string, args CommitContainerImageArgs) (err error) {

	url := fmt.Sprintf("%s/v3/containers/%s/commit", p.host, ip)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// POST /v3/containers/<ip>/exec
func (p *qcosClientImp) ExecContainer(
	ctx context.Context, ip string, args ExecContainerArgs) (
	ret ExecContainerRet, err error) {

	url := fmt.Sprintf("%s/v3/containers/%s/exec", p.host, ip)
	err = p.client.CallWithJson(ctx, &ret, "POST", url, args)
	return
}

// POST /v3/containers/<ip>/exec/<execId>/resize
func (p *qcosClientImp) ResizeContainerExecTerm(ctx context.Context,
	ip string, execID string, args ResizeContainerExecTermArgs) (err error) {

	url := fmt.Sprintf("%s/v3/containers/%s/exec/%s/resize", p.host, ip, execID)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// POST /v3/containers/<ip>/exec/<execId>/start
func (p *qcosClientImp) StartContainerExec(ctx context.Context, ip string, execID string, args StartContainerExecArgs, opts StartContainerExecOpts) (err error) {

	url := fmt.Sprintf("%s/v3/containers/%s/exec/%s/start", p.host, ip, execID)
	defer func() {
		if opts.ErrorCh != nil {
			opts.ErrorCh <- err
		}
	}()

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "tcp")

	if p.kmac != nil {
		if err = p.kmac.SignRequest(req); err != nil {
			return
		}
	}

	isSSL := false
	host := req.Host
	if req.URL.Scheme == "https" {
		isSSL = true
	}
	if !strings.Contains(host, ":") {
		if isSSL {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	var conn net.Conn
	if !isSSL {
		conn, err = net.DialTimeout("tcp", host, 10*time.Second)
		if err != nil {
			return fmt.Errorf("net.DialTimeout err: %v", err)
		}
		defer conn.Close()
	} else {
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", host, nil)
		if err != nil {
			return fmt.Errorf("tls.DialWithDialer err: %v", err)
		}
		defer conn.Close()
	}

	err = req.Write(conn)
	if err != nil {
		return fmt.Errorf("try send request: %v", err)
	}

	buf := bufio.NewReader(conn)
	resp, err := http.ReadResponse(buf, req)
	if err != nil {
		return fmt.Errorf("try receive response: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("try read response body: %v %v", string(b), err)
	}

	if resp.StatusCode != http.StatusSwitchingProtocols || !isUpgradeTCP(resp.Header) {
		return fmt.Errorf("not a upgrade proto code: %d %s", resp.StatusCode, string(b))
	}

	if opts.ReadyCh != nil {
		opts.ReadyCh <- struct{}{}
		<-opts.ReadyCh
	}

	errch := make(chan error, 2)
	go func() {
		_, err := stdCopy(opts.OutStream, opts.ErrStream, conn)
		errch <- err
	}()
	go func() {
		_, err := io.Copy(conn, opts.InStream)
		errch <- err
	}()

	err = <-errch
	if err == io.ErrClosedPipe {
		err = nil
	}
	if err != nil {
		err = fmt.Errorf("copy data: %v", err)
	}
	return nil
}

func isUpgradeTCP(headers http.Header) bool {
	return strings.Contains(strings.ToLower(headers.Get("Connection")), "upgrade") && headers.Get("Upgrade") != ""
}

// PUT /v3/containers/<ip>/webdav/files/<filePath>
func (p *qcosClientImp) UploadToContainer(
	ctx context.Context, ip string, filePath string, rd io.Reader) (
	err error) {

	url := fmt.Sprintf(
		path.Join("%s/v3/containers/%s/webdav/files/", filePath), p.host, ip)

	p.logger.WithField("url", url).Debug("webdav request")

	req, err := http.NewRequest("PUT", url, rd)
	if err != nil {
		return
	}
	resp, err := p.client.Do(ctx, req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	p.logger.WithField("body", string(text)).WithField("code",
		resp.Status).WithField("header", resp.Header).Debug("webdav result")

	switch resp.StatusCode {
	case http.StatusCreated:
		// do nothing
	case http.StatusNotFound:
		err = ErrNoSuchEntry
	default:
		err = ErrResultError
	}

	return
}

// GET /v3/containers/<ip>/webdav/files/<filePath>
func (p *qcosClientImp) DownloadFromContainer(
	ctx context.Context, ip string, filePath string) (
	rc io.ReadCloser, err error) {

	url := fmt.Sprintf(
		path.Join("%s/v3/containers/%s/webdav/files/", filePath), p.host, ip)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	p.logger.WithField("url", url).WithField("method", req.Method).Debug("webdav request")

	resp, err := p.client.Do(ctx, req)
	if err != nil {
		return
	}
	rc = resp.Body
	p.logger.WithField("code",
		resp.Status).WithField("header", resp.Header).Debug("webdav result")

	switch resp.StatusCode {
	case http.StatusOK:
		// do nothing
	case http.StatusNotFound:
		err = ErrNoSuchEntry
		rc.Close()
	default:
		err = ErrResultError
		rc.Close()
	}

	return
}

// PROPFIND /v3/containers/<ip>/webdav/files/<filePath>
func (p *qcosClientImp) StatContainerFile(ctx context.Context, ip string,
	filePath string, args StatContainerFileArgs) (rc io.ReadCloser, err error) {

	url := fmt.Sprintf(
		path.Join("%s/v3/containers/%s/webdav/files/", filePath), p.host, ip)

	req, err := http.NewRequest("PROPFIND", url, nil)
	if err != nil {
		return
	}
	if args.Depth == -1 {
		req.Header.Add("Depth", "infinity")
	} else {
		req.Header.Add("Depth", strconv.Itoa(args.Depth))
	}
	p.logger.WithField("url", url).WithField("method", req.Method).Debug("webdav request")

	resp, err := p.client.Do(ctx, req)
	if err != nil {
		return
	}
	rc = resp.Body
	p.logger.WithField("code",
		resp.Status).WithField("header", resp.Header).Debug("webdav result")

	switch resp.StatusCode {
	case http.StatusOK, MultiStatus:
		// do nothing
	case http.StatusMethodNotAllowed, http.StatusNotFound:
		err = ErrNoSuchEntry
		rc.Close()
	default:
		err = ErrResultError
		rc.Close()
	}

	return

}

// MKCOL /v3/containers/<ip>/webdav/files/<filePath>
func (p *qcosClientImp) MkdirInContainer(ctx context.Context,
	ip string, filePath string) (err error) {

	url := fmt.Sprintf(
		path.Join("%s/v3/containers/%s/webdav/files/", filePath), p.host, ip)

	req, err := http.NewRequest("MKCOL", url, nil)
	if err != nil {
		return
	}
	p.logger.WithField("url", url).WithField("method", req.Method).Debug("webdav request")

	resp, err := p.client.Do(ctx, req)
	if err != nil {
		return
	}
	rc := resp.Body
	p.logger.WithField("code",
		resp.Status).WithField("header", resp.Header).Debug("webdav result")

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		// do nothing
	case http.StatusNotFound:
		err = ErrNoSuchEntry
		rc.Close()
	default:
		err = ErrResultError
		rc.Close()
	}

	return

}

// GET /v3/logs/containers/<ip>/realtime?since=<since>&tail=<tail>
func (p *qcosClientImp) GetContainerLogsRealtime(ctx context.Context, ip, since, tail string, opts GetContainerLogsRealtimeOpts) (stream io.ReadCloser, err error) {
	url := fmt.Sprintf("%s/v3/logs/containers/%s/realtime?since=%s&tail=%s", p.host, ip, since, tail)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "tcp")

	if p.kmac != nil {
		if err = p.kmac.SignRequest(req); err != nil {
			return
		}
	}

	isSSL := false
	host := req.Host
	if req.URL.Scheme == "https" {
		isSSL = true
	}
	if !strings.Contains(host, ":") {
		if isSSL {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	var conn net.Conn
	if !isSSL {
		conn, err = net.DialTimeout("tcp", host, 10*time.Second)
		if err != nil {
			conn.Close()
			err = fmt.Errorf("net.DialTimeout err: %v", err)
			return
		}
	} else {
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", host, nil)
		if err != nil {
			conn.Close()
			err = fmt.Errorf("tls.DialWithDialer err: %v", err)
			return
		}
	}

	err = req.Write(conn)
	if err != nil {
		err = fmt.Errorf("try send request: %v", err)
		return
	}

	buf := bufio.NewReader(conn)
	resp, err := http.ReadResponse(buf, req)
	if err != nil {
		err = fmt.Errorf("try receive response: %v", err)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("try read response body: %v %v", string(b), err)
		return
	}

	if resp.StatusCode != http.StatusSwitchingProtocols || !isUpgradeTCP(resp.Header) {
		err = fmt.Errorf("get container logs realtime failed: %d %s", resp.StatusCode, string(b))
		return
	}

	stream, out := io.Pipe()

	go func() {

		defer func() {
			resp.Body.Close()
			conn.Close()
			stream.Close()
			out.Close()
		}()

		errCh := make(chan error, 1)
		go func() {
			_, err := io.Copy(out, conn)
			if err == syscall.EPIPE || err == io.ErrClosedPipe || err == io.EOF {
				err = nil
			} else {
				err = fmt.Errorf("copy data from conn: %v", err)
			}

			errCh <- err
		}()

		select {
		case e := <-errCh:
			opts.ErrorCh <- e
		case <-opts.ExitCh:
		}
	}()

	return
}

// GET /v3/logs/search/<repoType>?q=<queryString>&from=<from>&size=<size>&sort=<sort>&timeout=<timeout>
func (p *qcosClientImp) SearchContainerLogs(ctx context.Context, args SearchContainerLogsArgs) (res LogsSearchResult, err error) {
	if args.RepoType == "" {
		err = fmt.Errorf("RepoType could not be empty")
		return
	}

	queryURL := fmt.Sprintf("%s/v3/logs/search/%s", p.host, args.RepoType)
	params := make([]string, 0)
	if args.Query != "" {
		params = append(params, fmt.Sprintf("q=%s", url.QueryEscape(args.Query)))
	}
	if args.From != 0 {
		params = append(params, fmt.Sprintf("from=%d", args.From))
	}
	if args.Size != 0 {
		params = append(params, fmt.Sprintf("size=%d", args.Size))
	}
	if args.Sort != "" {
		params = append(params, fmt.Sprintf("sort=%s", args.Sort))
	}
	if len(params) > 0 {
		queryURL += "?" + strings.Join(params, "&")
	}

	err = p.client.Call(ctx, &res, "GET", queryURL)
	return
}

// generated from /Users/song/qbox/product/qcos/qcc-apidocs/includes/_aps.md
// Cient

// GET /v3/aps
func (p *qcosClientImp) ListAps(
	ctx context.Context, args ListApsArgs) (ret []ListApInfo, err error) {

	var query string
	if args.Service != "" {
		//QCOSD API框架无法正确处理Query中?和=的转义，对?和=转义会无法路由到正确的处理函数，返回非预期的结果。
		// = 之后的字符串需要做转义，因为空格无法正确处理，但title中空格是合法字符。
		query = "?service=" + url.QueryEscape(fmt.Sprintf("%s", args.Service))

	}
	if args.Stack != "" {
		query = "?stack=" + url.QueryEscape(fmt.Sprintf("%s", args.Stack))
	}
	if args.Title != "" {
		query = "?title=" + url.QueryEscape(fmt.Sprintf("%s", args.Title))
	}
	url := fmt.Sprintf("%s/v3/aps%s", p.host, query)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

//GET /v3/aps?stack=<stack> | GET /v3/aps?service=<service> | GET /v3/aps?title=<title>
func (p *qcosClientImp) ListApsFilter(ctx context.Context, filterKey string, filterValue string) (ret []ListApInfo, err error) {
	url := fmt.Sprintf("%s/v3/aps?%s=%s", p.host, filterKey, filterValue)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// POST /v3/aps
func (p *qcosClientImp) CreateAp(ctx context.Context, args CreateApArgs) (ret ListApInfo, err error) {
	url := fmt.Sprintf("%s/v3/aps", p.host)
	err = p.client.CallWithJson(ctx, &ret, "POST", url, args)
	return
}

// GET  /v3/aps/search?ip=<IP> | /v3/aps/search?domain=<domain> | /v3/aps/search?host=<host>
// mode : ip | domain | host ;  searchArg: <IP> | <domain> | <host>
func (p *qcosClientImp) SearchAp(ctx context.Context, mode string, searchArg string) (ret FullApInfo, err error) {
	url := fmt.Sprintf("%s/v3/aps/search?%s=%s", p.host, mode, searchArg)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// GET  /v3/aps/<apid>
func (p *qcosClientImp) GetAp(ctx context.Context, apid string) (ret FullApInfo, err error) {
	url := fmt.Sprintf("%s/v3/aps/%s", p.host, apid)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// POST /v3/aps/<apid>
func (p *qcosClientImp) UpdateAp(ctx context.Context, apid string, args SetApDescArgs) (err error) {
	url := fmt.Sprintf("%s/v3/aps/%s", p.host, apid)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// POST /v3/aps/<apid>/<port>
func (p *qcosClientImp) SetApPort(ctx context.Context, apid string, port string, args SetApPortArgs) (err error) {
	url := fmt.Sprintf("%s/v3/aps/%s/%s", p.host, apid, port)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// DELETE /v3/aps/<apid>/<port>
func (p *qcosClientImp) DeleteApPort(ctx context.Context, apid string, port string) (err error) {
	url := fmt.Sprintf("%s/v3/aps/%s/%s", p.host, apid, port)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

// POST /v3/aps/<apid>/portrange/<from>/<to>
func (p *qcosClientImp) SetApPortRange(
	ctx context.Context, apid string, fromPort string, toPort string, args SetApPortRangeArgs) (err error) {
	url := fmt.Sprintf("%s/v3/aps/%s/portrange/%s/%s", p.host, apid, fromPort, toPort)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// DELETE /v3/aps/<apid>/portrange/<from>/<to>
func (p *qcosClientImp) DeleteApPortRange(
	ctx context.Context, apid string, fromPort string, toPort string) (err error) {
	url := fmt.Sprintf("%s/v3/aps/%s/portrange/%s/%s", p.host, apid, fromPort, toPort)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

// GET  /v3/aps/<apid>/<port>/healthcheck
func (p *qcosClientImp) GetHealthcheck(ctx context.Context, apid string, port string) (ret map[string]string, err error) {
	url := fmt.Sprintf("%s/v3/aps/%s/%s/healthcheck", p.host, apid, port)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// DELETE /v3/aps/<apid>
func (p *qcosClientImp) DeleteAp(ctx context.Context, apid string) (err error) {
	url := fmt.Sprintf("%s/v3/aps/%s", p.host, apid)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

// POST /v3/aps/<apid>/<port>/setcontainer
func (p *qcosClientImp) ApSetContainer(ctx context.Context, apid string, port string, args []SetApContainerOptionsArgs) (err error) {
	url := fmt.Sprintf("%s/v3/aps/%s/%s/setcontainer", p.host, apid, port)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// POST /v3/aps/<apid>/publish
func (p *qcosClientImp) PublishUserDomain(ctx context.Context, apid string, args SetUserDomainArgs) (err error) {
	url := fmt.Sprintf("%s/v3/aps/%s/publish", p.host, apid)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// POST /v3/aps/<apid>/unpublish
func (p *qcosClientImp) UnpublishUserDomain(ctx context.Context, apid string, args SetUserDomainArgs) (err error) {
	url := fmt.Sprintf("%s/v3/aps/%s/unpublish", p.host, apid)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// GET /v3/aps/providers
func (p *qcosClientImp) ListProviders(ctx context.Context) (ret []string, err error) {
	url := fmt.Sprintf("%s/v3/aps/providers", p.host)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// GET /v3/jobs
func (p *qcosClientImp) ListJobs(ctx context.Context) (ret []JobInfo, err error) {
	url := fmt.Sprintf("%s/v3/jobs", p.host)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// GET /v3/jobs/<name>
func (p *qcosClientImp) GetJob(ctx context.Context, name string) (ret JobInfo, err error) {
	url := fmt.Sprintf("%s/v3/jobs/%s", p.host, name)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// DELETE /v3/jobs/<name>
func (p *qcosClientImp) DeleteJob(ctx context.Context, name string) (err error) {
	url := fmt.Sprintf("%s/v3/jobs/%s", p.host, name)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

// POST /v3/jobs
func (p *qcosClientImp) CreateJob(ctx context.Context, args CreateJobArgs) (err error) {
	url := fmt.Sprintf("%s/v3/jobs", p.host)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// POST /v3/jobs/<name>
func (p *qcosClientImp) UpdateJob(ctx context.Context, name string, args UpdateJobArgs) (err error) {
	url := fmt.Sprintf("%s/v3/jobs/%s", p.host, name)
	err = p.client.Call(ctx, nil, "POST", url)
	return
}

// POST /v3/jobs/<name>/run
func (p *qcosClientImp) RunJob(ctx context.Context, name string, args RunJobArgs) (ret JobInstanceID, err error) {
	url := fmt.Sprintf("%s/v3/jobs/%s/run", p.host, name)
	err = p.client.CallWithJson(ctx, &ret, "POST", url, args)
	return
}

// GET /v3/jobs/<name>/instances/<id>
func (p *qcosClientImp) GetJobInstance(ctx context.Context, name string, id string) (ret JobInstance, err error) {
	url := fmt.Sprintf("%s/v3/jobs/%s/instances/%s", p.host, name, id)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// DELETE /v3/jobs/<name>/instances/<id>
func (p *qcosClientImp) DeleteJobInstance(ctx context.Context, name string, id string) (err error) {
	url := fmt.Sprintf("%s/v3/jobs/%s/instances/%s", p.host, name, id)
	err = p.client.Call(ctx, nil, "DELETE", url)
	return
}

// POST /v3/jobs/<name>/instances/<id>/stop
func (p *qcosClientImp) StopJobInstance(ctx context.Context, name string, id string) (err error) {
	url := fmt.Sprintf("%s/v3/jobs/%s/instances/%s/stop", p.host, name, id)
	err = p.client.Call(ctx, nil, "POST", url)
	return
}

//  POST /v3/alert/aps/<apid>
func (p *qcosClientImp) UpdateApAlert(ctx context.Context, apid string, args UpdateApAlertArgs) (err error) {
	url := fmt.Sprintf("%s/v3/alert/aps/%s", p.host, apid)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// DELETE /v3/alert/aps/<apid>
func (p *qcosClientImp) DeleteApAlert(ctx context.Context, apid string, level string) (err error) {
	url := fmt.Sprintf("%s/v3/alert/aps/%s", p.host, apid)
	err = p.client.CallWithJson(ctx, nil, "DELETE", url, AlertLevelArgs{Level: level})
	return
}

//  GET /v3/alert/aps/<apid>?level=<level>
func (p *qcosClientImp) GetApAlert(ctx context.Context, apid string, level string) (ret []ApAlertInfo, err error) {
	var query string
	if level != "" {
		query = fmt.Sprintf("?level=%s", level)
	}
	url := fmt.Sprintf("%s/v3/alert/aps/%s%s", p.host, apid, query)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// POST /v3/alert/stacks/<stackName>/services/<serviceName>
func (p *qcosClientImp) UpdateServiceAlert(ctx context.Context, stack, service string, args UpdateContainerAlertArgs) (err error) {
	url := fmt.Sprintf("%s/v3/alert/stacks/%s/services/%s", p.host, stack, service)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// POST /v3/alert/stacks/<stackName>/services/<serviceName>/all
func (p *qcosClientImp) UpdateAllContainerAlert(ctx context.Context, stack, service string, args UpdateContainerAlertArgs) (err error) {
	url := fmt.Sprintf("%s/v3/alert/stacks/%s/services/all", p.host, stack, service)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// POST /v3/alert/containers/<ip>
func (p *qcosClientImp) UpdateContainerAlert(ctx context.Context, ip string, args UpdateContainerAlertArgs) (err error) {
	url := fmt.Sprintf("%s/v3/alert/containers/%s", p.host, ip)
	err = p.client.CallWithJson(ctx, nil, "POST", url, args)
	return
}

// DELETE /v3/alert/stacks/<stackName>/services/<serviceName>
func (p *qcosClientImp) DeleteServiceAlert(ctx context.Context, stack, service string, level string) (err error) {
	url := fmt.Sprintf("%s/v3/alert/stacks/%s/services/%s", p.host, stack, service)
	err = p.client.CallWithJson(ctx, nil, "DELETE", url, AlertLevelArgs{Level: level})
	return
}

// DELETE /v3/alert/containers/<ip>
func (p *qcosClientImp) DeleteContainerAlert(ctx context.Context, ip string, level string) (err error) {
	url := fmt.Sprintf("%s/v3/alert/containers/%s", p.host, ip)
	err = p.client.CallWithJson(ctx, nil, "DELETE", url, AlertLevelArgs{Level: level})
	return
}

// GET /v3/alert/stacks/<stackName>/services/<serviceName>?level=<level>
func (p *qcosClientImp) GetServiceAlert(ctx context.Context, stack, service string, level string) (ret []ContainerAlertInfo, err error) {
	var query string
	if level != "" {
		query = fmt.Sprintf("?level=%s", level)
	}
	url := fmt.Sprintf("%s/v3/alert/stacks/%s/services/%s%s", p.host, stack, service, query)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

// GET /v3/alert/containers/<ip>?level=<level>
func (p *qcosClientImp) GetContainerAlert(ctx context.Context, ip string, level string) (ret []ContainerAlertInfo, err error) {
	var query string
	if level != "" {
		query = fmt.Sprintf("?level=%s", level)
	}
	url := fmt.Sprintf("%s/v3/alert/containers/%s%s", p.host, ip, query)
	err = p.client.Call(ctx, &ret, "GET", url)
	return
}

func (p *qcosClientImp) wait4StackRunning(stackName string, timeout time.Duration) (err error) {
	if stackName == "" {
		stackName = DefaultStack
	}

	done := make(chan struct{})

	go func() {
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		for {
			stackInfo, err := p.GetStack(ctx, stackName)
			if err != nil {
				if err.Error() == "context deadline exceeded" {
					break
				}
			} else if err == nil && stackInfo.IsDeployed == true && stackInfo.Status == "RUNNING" {
				runningNum := 0
				for _, svcName := range stackInfo.Services {
					err = p.wait4ServiceRunning(stackName, svcName, timeout)
					if err != nil {
						break
					}
					runningNum += 1
				}
				if runningNum == len(stackInfo.Services) {
					done <- struct{}{}
					break
				}
			}
			time.Sleep(time.Second)
		}
	}()

	select {
	case <-time.After(timeout):
		return errors.New("Timeout in wait4StackRunning")
	case <-done:
		return nil
	}
}

func (p *qcosClientImp) wait4ServiceRunning(stackName string, serviceName string, timeout time.Duration) (err error) {
	if stackName == "" {
		stackName = DefaultStack
	}

	done := make(chan struct{})

	go func() {
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		for {
			svcInfo, err := p.GetServiceInspect(ctx, stackName, serviceName)
			if err != nil {
				if err.Error() == "context deadline exceeded" {
					break
				}
			} else if err == nil && svcInfo.State == "DEPLOYED" && svcInfo.Status == "RUNNING" {
				runningNum := 0
				for _, contIp := range svcInfo.ContainerIPs {
					ctx, _ := context.WithTimeout(context.Background(), timeout)
					contInfo, err := p.GetContainerInspect(ctx, contIp)
					if err != nil || contInfo.Status != "RUNNING" {
						break
					}
					runningNum += 1
				}
				if runningNum == len(svcInfo.ContainerIPs) {
					done <- struct{}{}
					break
				}
			}
			time.Sleep(time.Second)
		}
	}()

	select {
	case <-time.After(timeout):
		return errors.New("Timeout in wait4ServiceRunning")
	case <-done:
		return nil
	}
}

func (p *qcosClientImp) wait4ServiceStopped(stackName string, serviceName string, timeout time.Duration) (err error) {
	if stackName == "" {
		stackName = DefaultStack
	}

	done := make(chan struct{})

	go func() {
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		for {
			svcInfo, err := p.GetServiceInspect(ctx, stackName, serviceName)
			if err != nil {
				if err.Error() == "context deadline exceeded" {
					break
				}
			} else if err == nil && svcInfo.State == "STOPPED" && svcInfo.Status == "NOT-RUNNING" {
				runningNum := 0
				for _, contIp := range svcInfo.ContainerIPs {
					ctx, _ := context.WithTimeout(context.Background(), timeout)
					contInfo, err := p.GetContainerInspect(ctx, contIp)
					if err != nil || contInfo.Status != "EXITED" {
						break
					}
					runningNum += 1
				}
				if runningNum == len(svcInfo.ContainerIPs) {
					done <- struct{}{}
					break
				}
			}
			time.Sleep(time.Second)
		}
	}()

	select {
	case <-time.After(timeout):
		return errors.New("Timeout in wait4ServiceRunning")
	case <-done:
		return nil
	}
}
