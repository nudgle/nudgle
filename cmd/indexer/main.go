package main

import (
	"github.com/nudgle/nudgle/internal/config"
	"github.com/nudgle/nudgle/internal/indexer"
)

func init() {
	config.GetFlags(config.IndexerService)
}

func main() {
	indexer.Start()
}
