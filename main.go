package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path"
	"time"

	"github.com/ashep/sbk/config"
	"github.com/ashep/sbk/notifier"
	"github.com/ashep/sbk/rdiffbackup"
	"github.com/ashep/sbk/telegram"
)

func main() {
	cfgPath := flag.String("c", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.ParseFromFile(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config file: %s", err)
		return
	}

	if cfg.LogDir == "" {
		cfg.LogDir = "./log"
	}
	if err := os.MkdirAll(cfg.LogDir, 0o755); err != nil {
		log.Fatalf("failed to create log directory: %s", err)
		return
	}

	time.Local = time.UTC

	ntf := notifier.NewNoop()
	if cfg.Notifications.Telegram != nil {
		ntf = notifier.NewTelegram(telegram.New(cfg.Notifications.Telegram.Token), cfg.Notifications.Telegram.ChatId)
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	if len(cfg.Files.Sources) != 0 {
		logFilename := path.Join(cfg.LogDir, time.Now().Format("20060102-150405-files.log"))
		rb := rdiffbackup.New(cfg.Files.Verbosity, ntf)
		rb.BatchBackup(ctx, cfg.Files.Sources, cfg.Files.Exclude, cfg.Files.Destination, logFilename)
	}

	ctxCancel()
}
