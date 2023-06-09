package util

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
)

func AbsPath(p string) string {
	if path.IsAbs(p) {
		return p
	}

	wd, err := os.Getwd()
	if err != nil {
		return p
	}

	return path.Join(wd, p)
}

func StreamCommand(ctx context.Context, outW, errW io.Writer, name string, args []string) error {
	c := exec.CommandContext(ctx, name, args...)

	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}

	if err = c.Start(); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go Stream(ctx, &wg, outW, errW, stdout, stderr)

	err = c.Wait()
	wg.Wait()

	return err
}

// Stream reads from inR, errR and writes read bytes to outW, errW correspondingly
func Stream(ctx context.Context, wg *sync.WaitGroup, outW, errW io.Writer, inR, errR io.Reader) {
	buf := make([]byte, 64<<10) // 64 Kb
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := flush(inR, outW, buf)
			if err != nil {
				if !(err == io.EOF || errors.Is(err, fs.ErrClosed)) {
					log.Print(err)
				}
				return
			}

			err = flush(errR, errW, buf)
			if err != nil {
				if !(err == io.EOF || errors.Is(err, fs.ErrClosed)) {
					log.Print(err)
				}
				return
			}

		}
	}
}

func flush(r io.Reader, w io.Writer, buf []byte) error {
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	_, err = w.Write(buf[:n])
	if err != nil {
		return err
	}

	return nil
}
