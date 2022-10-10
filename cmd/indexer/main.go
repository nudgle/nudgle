package main

import (
	"nudgle/internal/config"
	"nudgle/internal/indexer"
)

func init() {
	config.GetFlags(config.IndexerService)
}

func main() {
	indexer.Start()
}
