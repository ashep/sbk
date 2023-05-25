package util

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log"
	"os/exec"
	"sync"
)

func StreamCommand(ctx context.Context, out, errs io.Writer, name string, args []string) error {
	c := exec.CommandContext(ctx, name, args...)

	cmdStdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	cmdStderr, err := c.StderrPipe()
	if err != nil {
		return err
	}

	if err = c.Start(); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go Stream(ctx, &wg, out, errs, cmdStdout, cmdStderr)

	err = c.Wait()
	wg.Wait()

	return err
}

func Stream(ctx context.Context, wg *sync.WaitGroup, out, errsOut io.Writer, in, errsIn io.Reader) {
	buf := make([]byte, 64<<10) // 64 Kb
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := flush(in, out, buf)
			if err != nil {
				if !(err == io.EOF || errors.Is(err, fs.ErrClosed)) {
					log.Print(err)
				}
				return
			}

			err = flush(errsIn, errsOut, buf)
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
	if _, err := r.Read(buf); err != nil {
		return err
	}

	if _, err := w.Write(buf); err != nil {
		return err
	}

	return nil
}
