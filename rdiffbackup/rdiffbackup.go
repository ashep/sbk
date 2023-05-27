package rdiffbackup

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ashep/sbk/util"
)

type notifier interface {
	NotifySuccess(msg string) error
	NotifyError(msg string) error
}

type RdiffBackup struct {
	verbosity int
}

func New(verbosity int) *RdiffBackup {
	return &RdiffBackup{
		verbosity: verbosity,
	}
}

func (r *RdiffBackup) BackupMany(ctx context.Context, sources, exclude []string, target string, ntf notifier) {
	for _, src := range sources {
		dst := target + "/" + src

		err := r.Backup(ctx, src, dst, exclude)
		if err != nil && err == context.Canceled {
			return
		} else if err != nil {
			msg := fmt.Sprintf("backup failed: %s: %s", src, err)
			log.Print(msg)
			if ntf != nil {
				ntf.NotifyError(msg)
			}
			continue
		}

		msg := fmt.Sprintf("backup ok: %s", src)

		log.Print(msg)

		if ntf != nil {
			if err := ntf.NotifySuccess(msg); err != nil {
				log.Printf("failed to send notification: %s", err)
			}
		}
	}
}

func (r *RdiffBackup) Backup(ctx context.Context, src, dst string, exclude []string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	args := make([]string, 0)
	for _, exc := range exclude {
		args = append(args, "--exclude", exc)
	}

	args = append(args, "-v", strconv.Itoa(r.verbosity))
	args = append(args, src, dst)

	return util.StreamCommand(ctx, os.Stdout, os.Stderr, "rdiff-backup", args)
}
