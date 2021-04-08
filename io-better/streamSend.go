package io_better

import (
	"context"
	"io"

	go_group "bernieyangmh.com/TinyComponent/go-group"
)

func ParallelStreamIO(writer io.Writer, reader io.Reader, wf func(w io.Writer, r io.Reader) error, rf func(r io.Reader, w io.Writer) error) (err error) {
	var (
		rp, wp = io.Pipe()
		gGroup = &go_group.Group{}
	)

	gGroup.Go(func(ctx context.Context) error {
		if err = rf(reader, wp); err != nil {
			return wp.Close()
		}
		return nil
	})

	gGroup.Go(func(ctx context.Context) error {
		if err = wf(writer, rp); err != nil {
			return wp.Close()
		}
		return nil
	})
	return gGroup.Wait()
}
