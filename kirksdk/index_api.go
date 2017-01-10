package kirksdk

import (
	"strings"
	"time"

	"golang.org/x/net/context"
)

type IndexAuthClient interface {
	GetConfig() (ret IndexAuthConfig)
	RequestAuthToken(ctx context.Context, scopes []string) (AuthToken, error)
}

type IndexClient interface {
	GetConfig() (ret IndexConfig)
	ListRepo(ctx context.Context, username string) (repos []*Repo, err error)
	ListRepoTags(ctx context.Context, username, repo string) (tags []*Tag, err error)
	ListRepoTagsPage(ctx context.Context, username, repo string, start, size int) (tags []*Tag, err error)
	GetImageConfig(ctx context.Context, username, repo, reference string) (res *ImageConfig, err error)
	DeleteRepoTag(ctx context.Context, username, repo, reference string) error
	CreateTagFromRepo(ctx context.Context, username, repo, tag string, from *ImageSpec) (result *ImageSpec, err error)
}

type AuthToken struct {
	Token     string    `json:"token"`
	ExpiresIn int64     `json:"expires_in":`
	IssuedAt  time.Time `json:"issued_at":`
}

type ImageSpec struct {
	Username  string `json:"username"`
	Repo      string `json:"repo"`
	Reference string `json:"reference"`
}

type Repo struct {
	Name string `json:"name"`
}

type Tag struct {
	Name    string      `json:"name"`
	Created time.Time   `json:"created"`
	Detail  ImageConfig `json:"detail"`
}

type ImageConfig struct {
	Digest Digest `json:"digest"`
	// TODO
	// use Struct instead of map[string]interface{}
	Config          map[string]interface{} `json:"config"`
	ContainerConfig map[string]interface{} `json:"container_config"`
	Created         time.Time              `json:"created"`
	Size            int64                  `json:"size"`
}

type Digest string

func (p Digest) String() string {
	return string(p)
}

func (p Digest) ID() string {
	s := string(p)
	if strings.HasPrefix(s, "sha256:") {
		s = s[7:]
	}
	return s
}
