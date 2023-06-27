package traefik_query_append_url

import (
	"context"
	"fmt"
	"net/http"
)

// Config is the configuration for this plugin
type Config struct {
	QueryParamName string `json:"queryParamName"`
}

// CreateConfig creates a new configuration for this plugin
func CreateConfig() *Config {
	return &Config{}
}

// QueryModification represents the basic properties of this plugin
type QueryModification struct {
	next   http.Handler
	name   string
	config *Config
}

// New creates a new instance of this plugin
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.QueryParamName == "" {
		config.QueryParamName = "url"
	}
	return &QueryModification{
		next:   next,
		name:   name,
		config: config,
	}, nil
}

func (q *QueryModification) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" || req.Method == "" {
		qry := req.URL.Query()
		url := fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.Path)
		qry.Add(q.config.QueryParamName, url)
		req.URL.RawQuery = qry.Encode()
		req.RequestURI = req.URL.RequestURI()
		q.next.ServeHTTP(rw, req)
	}
}
