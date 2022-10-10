package filter

import (
	"fmt"
	monitor "github.com/nudgle/nudgle/pkg/monitor/config"
	"sync"
)

var pluginManager Manager

func init() {
	once := sync.Once{}
	once.Do(func() {
		GetPluginManager()
	})
}

type (
	// ManagerImpl is a manager for the plugins. It stores and maps the plugins to the command.
	ManagerImpl struct {
		plugins sync.Map
	}

	Manager interface {
		HasPlugin(name string) bool
		AddPlugin(name string, plugin Handler)
		GetPlugin(name string) (Handler, error)
		Chain(config *monitor.MonitorServiceConfiguration, sender *chan string, tx string)
	}

	Handler interface {
		Execute(config *monitor.MonitorServiceConfiguration, tx string) string
	}
)

// Register a cmd to the manager.
func Register(name string, plugin Handler) {
	GetPluginManager().AddPlugin(name, plugin)
}

// GetPluginManager gets the cmd manager instance (singleton).
func GetPluginManager() Manager {
	if pluginManager == nil {
		pluginManager = &ManagerImpl{plugins: sync.Map{}}
	}
	return pluginManager
}

// HasPlugin Check if the cmd exists in the manager.
func (pm *ManagerImpl) HasPlugin(name string) bool {
	_, exists := pm.plugins.Load(name)
	return exists
}

// AddPlugin Add a cmd to the manager.
func (pm *ManagerImpl) AddPlugin(name string, plugin Handler) {
	if !pm.HasPlugin(name) {
		pm.plugins.Store(name, plugin)
	}
}

// GetPlugin returns the cmd, if found.
func (pm *ManagerImpl) GetPlugin(name string) (Handler, error) {
	mPlugin, exists := pm.plugins.Load(name)
	if exists {
		return mPlugin.(Handler), nil
	}

	return nil, fmt.Errorf("cmd doesnt exist")
}

func (pm *ManagerImpl) Chain(
	config *monitor.MonitorServiceConfiguration,
	sender *chan string,
	tx string,
) {
	pm.plugins.Range(func(key interface{}, value interface{}) bool {
		filterMap, ok := pm.plugins.Load(key)
		if !ok {
			return true
		}
		filter := filterMap.(Handler)
		output := filter.Execute(config, tx)
		if output != "" {
			*sender <- output
		}
		return true
	})
}
