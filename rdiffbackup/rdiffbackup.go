package rdiffbackup

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ashep/sbk/icon"
	"github.com/ashep/sbk/notifier"
	"github.com/ashep/sbk/util"
)

type RdiffBackup struct {
	verbosity int
	notifier  notifier.Notifier
}

func New(verbosity int, notifier notifier.Notifier) *RdiffBackup {
	return &RdiffBackup{
		verbosity: verbosity,
		notifier:  notifier,
	}
}

func (r *RdiffBackup) BatchBackup(ctx context.Context, sources, exclude []string, target, logPath string) {
	if len(sources) == 0 {
		return
	}

	reportOk := ""
	reportErr := ""

	for _, src := range sources {
		dst := target + "/" + src
		now := time.Now()

		if err := r.Backup(ctx, src, dst, logPath, exclude); err != nil {
			log.Printf("Backup failed: src: %s; dst: %s; err: %s", src, dst, err)

			if reportErr != "" {
				reportErr += "\n\n"
			}
			reportErr += icon.Error + " Backup failed\n\n"
			reportErr += "• *source:* `" + src + "`\n"
			reportErr += "• *target:* `" + dst + "`\n"
			reportErr += "• *error:* `" + err.Error() + "`\n"
			reportErr += "• *time:* `" + time.Since(now).String() + "`"

			if errors.Is(err, context.Canceled) {
				break
			}

			continue
		}

		log.Printf("Backup succeed: src: %s; dst: %s", src, dst)

		if reportOk != "" {
			reportOk += "\n\n"
		}
		reportOk += icon.Success + " Backup succeed\n\n"
		reportOk += "• *source:* `" + src + "`\n"
		reportOk += "• *target:* `" + dst + "`\n"
		reportOk += "• *time:* `" + time.Since(now).String() + "`"
	}

	report := ""
	if reportErr != "" {
		report += reportErr
	}
	if reportOk != "" {
		report += reportOk
	}

	if err := r.notifier.Notify(report); err != nil {
		log.Printf("failed to send report: %s", err)
	}
}

func (r *RdiffBackup) Backup(ctx context.Context, src, dst, logPath string, exclude []string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	logF, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		_ = logF.Close()
	}()

	args := make([]string, 0)
	for _, exc := range exclude {
		args = append(args, "--exclude", exc)
	}

	args = append(args, "-v", strconv.Itoa(r.verbosity))
	args = append(args, src, dst)

	return util.StreamCommand(ctx, logF, logF, "rdiff-backup", args)
}
