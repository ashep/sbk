package mysql

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/ashep/sbk/config"
	"github.com/ashep/sbk/icon"
	"github.com/ashep/sbk/notifier"
	"github.com/ashep/sbk/util"
)

type MySQL struct {
	ntf notifier.Notifier
}

func New(notifier notifier.Notifier) *MySQL {
	return &MySQL{ntf: notifier}
}

func (m *MySQL) BatchBackup(ctx context.Context, sources []config.DBSource, target, logPath string) {
	if len(sources) == 0 {
		return
	}

	reportOk := ""
	reportErr := ""
	host, _ := os.Hostname()

	for _, src := range sources {
		select {
		case <-ctx.Done():
			return
		default:
		}

		dst := util.AbsPath(target + "/" + src.Database + ".sql")
		now := time.Now()
		srcStr := fmt.Sprintf("%s:%d/%s", src.Host, src.Port, src.Database)
		logMsg := fmt.Sprintf("src: %s; dst: %s", srcStr, dst)

		log.Print("MySQL backup started: " + logMsg)

		if err := m.Backup(ctx, src, dst, logPath); err != nil {
			log.Printf("MySQL backup failed: %s; err: %s", logMsg, err)

			if reportErr != "" {
				reportErr += "\n\n"
			}
			reportErr += icon.Error + " MySQL backup failed\n\n"
			reportErr += "• *host:* `" + host + "`\n"
			reportErr += "• *source:* `" + srcStr + "`\n"
			reportErr += "• *target:* `" + dst + ".gz`\n"
			reportErr += "• *time:* `" + time.Since(now).String() + "`\n"
			reportErr += "• *error:* `" + err.Error() + "`\n"
			reportErr += "• *log:* `" + logPath + "`\n"

			if errors.Is(err, context.Canceled) {
				break
			}

			continue
		}

		log.Print("MySQL backup succeed: " + logMsg)

		if reportOk != "" {
			reportOk += "\n\n"
		}
		reportOk += icon.Success + " MySQL backup succeed\n\n"
		reportOk += "• *host:* `" + host + "`\n"
		reportOk += "• *source:* `" + srcStr + "`\n"
		reportOk += "• *target:* `" + dst + ".gz`\n"
		reportOk += "• *time:* `" + time.Since(now).String() + "`\n"
		reportOk += "• *log:* `" + logPath + "`\n"
	}

	report := ""
	if reportErr != "" {
		report += reportErr
	}
	if reportOk != "" {
		report += reportOk
	}

	if err := m.ntf.Notify(report); err != nil {
		log.Printf("failed to send report: %s", err)
	}
}

func (m *MySQL) Backup(ctx context.Context, src config.DBSource, dst, logPath string) error {
	dstDir := path.Dir(dst)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return err
	}

	outF, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		_ = outF.Close()
	}()

	args := make([]string, 0)
	args = append(args, "-h", src.Host)
	args = append(args, "-P", strconv.Itoa(src.Port))
	args = append(args, "-u", src.User)
	args = append(args, "-p"+src.Password)
	args = append(args, "--protocol", "tcp")
	args = append(args, "--log-error", logPath)
	args = append(args, "--tz-utc")
	args = append(args, "--skip-lock-tables")
	args = append(args, src.Database)

	mdErr := util.StreamCommand(ctx, outF, io.Discard, "mysqldump", args)
	logStat, err := os.Stat(logPath)
	if mdErr != nil {
		return mdErr
	}

	if err == nil && logStat.Size() == 0 {
		_ = os.Remove(logPath)
	}

	// TODO: don't discard stderr, write it to logPath
	return util.StreamCommand(ctx, io.Discard, io.Discard, "gzip", []string{"-9", "-f", dst})
}
