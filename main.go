package main

import (
	"context"
	"flag"
	"log"

	"github.com/ashep/sbk/config"
	"github.com/ashep/sbk/rdiffbackup"
)

func main() {
	cfgPath := flag.String("c", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.ParseFromFile(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config file: %s", err)
		return
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	rb := rdiffbackup.New(cfg.Files.Verbosity)
	rb.Backup(ctx, cfg.Files.Sources, cfg.Files.Exclude, cfg.Files.Destination)

	ctxCancel()
}
