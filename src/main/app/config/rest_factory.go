package config

import (
	"fmt"
	"regexp"
	"time"

	"github.com/arielsrv/golang-toolkit/rest"
	"github.com/go-chassis/go-archaius"
)

const DefaultPoolSize = 10
const DefaultPoolTimeout = 200
const DefaultSocketTimeout = 500

func ProvideRestClients() *RESTClientFactory {
	restPoolPattern := regexp.MustCompile(`rest\.pool\.([-_\w]+)\..+`)

	restPoolFactory := &RESTPoolFactory{builders: map[string]*rest.RequestBuilder{}}
	poolsNames := getNamesInKeys(restPoolPattern)
	for _, name := range poolsNames {
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

	restClientPattern := regexp.MustCompile(`rest\.client\.([-_\w]+)\..+`)

	restClientFactory := RESTClientFactory{clients: map[string]*rest.RequestBuilder{}}
	clientsNames := getNamesInKeys(restClientPattern)
	for _, name := range clientsNames {
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

func Contains(elements []string, target string) bool {
	for _, element := range elements {
		if target == element {
			return true
		}
	}
	return false
}

func getNamesInKeys(regex *regexp.Regexp) []string {
	var result []string
	configs := archaius.GetConfigs()
	for key := range configs {
		match := regex.FindStringSubmatch(key)
		for i := range regex.SubexpNames() {
			if i > 0 && i <= len(match) {
				if !Contains(result, match[1]) {
					result = append(result, match[1])
				}
			}
		}
	}
	return result
}
