package config

import (
	"fmt"
	"regexp"
	"time"

	"github.com/arielsrv/ikp_go-restclient/rest"
	"github.com/arielsrv/pets-api/src/main/app/helpers/arrays"
	"github.com/go-chassis/go-archaius"
)

const (
	DefaultPoolSize      = 10
	DefaultPoolTimeout   = 200
	DefaultSocketTimeout = 500
)

const (
	RestPoolPattern   = `rest\.pool\.([-_\w]+)\..+`
	RestClientPattern = `rest\.client\.([-_\w]+)\..+`
)

var restPoolPattern = regexp.MustCompile(RestPoolPattern)
var restClientPattern = regexp.MustCompile(RestClientPattern)

func ProvideRestClients() *RESTClientFactory {
	restPoolFactory := &RESTPoolFactory{builders: map[string]*rest.RequestBuilder{}}
	poolNames := getNamesInKeys(restPoolPattern)
	for _, name := range poolNames {
		restPool := &rest.RequestBuilder{
			Timeout:        time.Millisecond * time.Duration(TryInt(fmt.Sprintf("rest.pool.%s.pool.timeout", name), DefaultPoolTimeout)),
			ConnectTimeout: time.Millisecond * time.Duration(TryInt(fmt.Sprintf("rest.pool.%s.pool.connection-timeout", name), DefaultSocketTimeout)),
			CustomPool: &rest.CustomPool{
				MaxIdleConnsPerHost: TryInt(fmt.Sprintf("rest.pool.%s.pool.size", name), DefaultPoolSize),
			},
		}
		restPoolFactory.Add(restPool)
		restPoolFactory.Register(name, restPool)
	}

	restClientFactory := RESTClientFactory{clients: map[string]*rest.RequestBuilder{}}
	clientNames := getNamesInKeys(restClientPattern)
	for _, name := range clientNames {
		poolName := String(fmt.Sprintf("rest.client.%s.pool", name))
		pool := restPoolFactory.GetPool(poolName)
		restClientFactory.Register(name, pool)
	}

	return &restClientFactory
}

type RESTPoolFactory struct {
	restPools []*rest.RequestBuilder
	builders  map[string]*rest.RequestBuilder
}

func (r *RESTPoolFactory) Add(rb *rest.RequestBuilder) {
	r.restPools = append(r.restPools, rb)
}

func (r *RESTPoolFactory) Register(name string, rb *rest.RequestBuilder) {
	r.builders[name] = rb
}

func (r *RESTPoolFactory) GetPool(name string) *rest.RequestBuilder {
	return r.builders[name]
}

type RESTClient struct {
	rest.RequestBuilder
}

type RESTClientFactory struct {
	clients map[string]*rest.RequestBuilder
}

func (r *RESTClientFactory) Register(name string, restPool *rest.RequestBuilder) {
	r.clients[name] = restPool
}

func (r *RESTClientFactory) Get(name string) *rest.RequestBuilder {
	return r.clients[name]
}

func (r *RESTClientFactory) GetClients() map[string]*rest.RequestBuilder {
	return r.clients
}

func getNamesInKeys(regex *regexp.Regexp) []string {
	var names []string
	configs := archaius.GetConfigs()
	for key := range configs {
		match := regex.FindStringSubmatch(key)
		for i := range regex.SubexpNames() {
			if i > 0 && i <= len(match) {
				group := match[1]
				if !arrays.Contains(names, group) {
					names = append(names, group)
				}
			}
		}
	}
	return names
}
