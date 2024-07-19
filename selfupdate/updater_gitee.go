package selfupdate

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"gitee.com/naibahq/go-gitee/gitee"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/peterbourgon/diskv"
)

// GiteeUpdater contains Gitee client and its context.
type GiteeUpdater struct {
	api       *gitee.APIClient
	apiCtx    context.Context
	validator Validator
	filters   []*regexp.Regexp
}

// NewGiteeUpdater creates a new updater instance. It initializes Gitee API client.
// If you set your API token to $GITEE_TOKEN, the client will use it.
func NewGiteeUpdater(config Config) (*GiteeUpdater, error) {
	token := config.APIToken
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	ctx := context.Background()

	filtersRe := make([]*regexp.Regexp, 0, len(config.Filters))
	for _, filter := range config.Filters {
		re, err := regexp.Compile(filter)
		if err != nil {
			return nil, fmt.Errorf("Could not compile regular expression %q for filtering releases: %v", filter, err)
		}
		filtersRe = append(filtersRe, re)
	}

	conf := gitee.NewConfiguration()
	cache := diskcache.NewWithDiskv(diskv.New(diskv.Options{}))
	transport := httpcache.NewTransport(cache)
	conf.HTTPClient = newHTTPClient(transport, token)

	client := gitee.NewAPIClient(conf)
	return &GiteeUpdater{api: client, apiCtx: ctx, validator: config.Validator, filters: filtersRe}, nil
}

// DefaultGiteeUpdater creates a new updater instance with default configuration.
// It initializes Gitee API client with default API base URL.
// If you set your API token to $GITEE_TOKEN, the client will use it.
func DefaultGiteeUpdater() *GiteeUpdater {
	token := os.Getenv("GITEE_TOKEN")
	// if token == "" {
	// 	token, _ = gitconfig.GithubToken()
	// }
	ctx := context.Background()
	conf := gitee.NewConfiguration()
	cache := diskcache.NewWithDiskv(diskv.New(diskv.Options{}))
	transport := httpcache.NewTransport(cache)
	conf.HTTPClient = newHTTPClient(transport, token)
	return &GiteeUpdater{api: gitee.NewAPIClient(conf), apiCtx: ctx}
}
