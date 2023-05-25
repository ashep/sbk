package rdiffbackup

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/ashep/sbk/util"
)

type RdiffBackup struct {
	verbosity int
}

func New(verbosity int) *RdiffBackup {
	return &RdiffBackup{
		verbosity: verbosity,
	}
}

func (r *RdiffBackup) Backup(ctx context.Context, sources, exclude []string, target string) {
	for _, src := range sources {
		dst := target + "/" + src

		if err := os.MkdirAll(dst, 0o755); err != nil {
			log.Print(err)
			continue
		}

		args := make([]string, 0)
		for _, exc := range exclude {
			args = append(args, "--exclude", exc)
		}

		args = append(args, "-v", strconv.Itoa(r.verbosity))
		args = append(args, src, dst)

		if err := util.StreamCommand(ctx, os.Stdout, os.Stderr, "rdiff-backup", args); err != nil {
			log.Print(err)
		}
	}
}
