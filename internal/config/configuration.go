package config

import (
	"flag"
	"fmt"
	"github.com/kkyr/fig"
	"github.com/nudgle/nudgle/internal/cache"
	indexer "github.com/nudgle/nudgle/pkg/indexer/config"
	monitor "github.com/nudgle/nudgle/pkg/monitor/config"
	cache2 "github.com/patrickmn/go-cache"
	"log"
	"os"
	"path/filepath"
)

const (
	IndexerConfigCacheKey string = "indexerConfig"
	MonitorConfigCacheKey string = "monitorConfig"
	IndexerService               = Service("indexer")
	MonitorService               = Service("monitor")
)

type Service string

func GetFlags(service Service) string {
	var (
		cacheKey            = ""
		memory              = cache.Memory()
		configFileName      *string
		defaultFileName     = "config"
		workingDirectory, _ = os.Getwd()
		configurationPath   = fmt.Sprintf("%s", workingDirectory)
	)
	switch service {
	case IndexerService:
		cacheKey = string(IndexerService)
		defaultFileName = string(IndexerService)
		break
	case MonitorService:
		cacheKey = string(MonitorService)
		defaultFileName = string(MonitorService)
		break
	}
	fileName := fmt.Sprintf("%s/config/%s.yaml", configurationPath, defaultFileName)
	configFileName = flag.String("config-file", fileName, "Path of the configuration file")
	flag.Parse()
	memory.Set(cacheKey, *configFileName, cache2.NoExpiration)

	return *configFileName
}

func GetServiceConfiguration(config interface{}, configFilePath string) {
	var (
		err error
	)

	err = fig.Load(config,
		fig.UseEnv("CUSTOMER_BAG"),
		fig.File(filepath.Base(configFilePath)),
		fig.Dirs(filepath.Dir(configFilePath), "../../config"),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func GetMonitorConfig() *monitor.MonitorServiceConfiguration {
	var (
		config              monitor.MonitorServiceConfiguration
		memory              = cache.Memory()
		isFound             bool
		cachedConfiguration interface{}
	)

	// Check if the configuration is cached
	cachedConfiguration, isFound = memory.Get(MonitorConfigCacheKey)
	if isFound {
		return cachedConfiguration.(*monitor.MonitorServiceConfiguration)
	}

	configFilePath, isFound := memory.Get(string(MonitorService))
	if !isFound {
		log.Fatal("Flags cache not found")
	}
	// Get configuration
	GetServiceConfiguration(&config, configFilePath.(string))

	memory.Set(MonitorConfigCacheKey, &config, cache2.NoExpiration)

	return &config
}

func GetIndexerConfig() *indexer.IndexerServiceConfiguration {
	var (
		config              indexer.IndexerServiceConfiguration
		memory              = cache.Memory()
		isFound             bool
		cachedConfiguration interface{}
	)

	// Check if the configuration is cached
	cachedConfiguration, isFound = memory.Get(IndexerConfigCacheKey)
	if isFound {
		return cachedConfiguration.(*indexer.IndexerServiceConfiguration)
	}

	configFilePath, isFound := memory.Get(string(IndexerService))
	if !isFound {
		log.Fatal("Flags cache not found")
	}
	// Get configuration
	GetServiceConfiguration(&config, configFilePath.(string))

	memory.Set(IndexerConfigCacheKey, &config, cache2.NoExpiration)

	return &config
}
