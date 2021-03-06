package remote

import (
	"fmt"

	"github.com/go-mesh/openlogging"
)

var configClientPlugins = make(map[string]func(options Options) (Client, error))

//const
const (
	LabelService     = "serviceName"
	LabelVersion     = "version"
	LabelEnvironment = "environment"
	LabelApp         = "app"
)

//DefaultClient is config server's client
var DefaultClient Client

//InstallConfigClientPlugin install a config client plugin
func InstallConfigClientPlugin(name string, f func(options Options) (Client, error)) {
	configClientPlugins[name] = f
	openlogging.GetLogger().Infof("Installed %s Plugin", name)
}

//Client is the interface of config server client, it has basic func to interact with config server
type Client interface {
	//PullConfigs pull all configs from remote
	PullConfigs(labels ...map[string]string) (map[string]interface{}, error)
	//PullConfig pull one config from remote
	PullConfig(key, contentType string, labels map[string]string) (interface{}, error)
	// PushConfigs push config to c
	PushConfigs(data map[string]interface{}, labels map[string]string) (map[string]interface{}, error)
	// DeleteConfigsByKeys delete config for c by keys
	DeleteConfigsByKeys(keys []string, labels map[string]string) (map[string]interface{}, error)
	//Watch get kv change results, you can compare them with local kv cache and refresh local cache
	Watch(f func(map[string]interface{}), errHandler func(err error), labels map[string]string) error
	Options() Options
}

//NewClient create config client implementation
func NewClient(name string, options Options) (Client, error) {
	plugins := configClientPlugins[name]
	if plugins == nil {
		return nil, fmt.Errorf("plugin [%s] not found", name)
	}
	DefaultClient, err := plugins(options)
	if err != nil {
		return nil, err
	}
	openlogging.GetLogger().Infof("%s plugin is enabled", name)
	return DefaultClient, nil
}
