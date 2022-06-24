// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntup

import (
	"fmt"
	"io"
)

type rbuff struct {
	p []byte // buffer of data to read from
	c int    // current position in buffer of data
}

func (r *rbuff) Read(p []byte) (int, error) {
	if r.c >= len(r.p) {
		return 0, io.EOF
	}
	n := copy(p, r.p[r.c:])
	r.c += n
	return n, nil
}

func (r *rbuff) ReadByte() (byte, error) {
	if r.c >= len(r.p) {
		return 0, io.EOF
	}
	v := r.p[r.c]
	r.c++
	return v, nil
}

func (r *rbuff) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.c = int(offset)
	case io.SeekCurrent:
		r.c += int(offset)
	case io.SeekEnd:
		r.c = len(r.p) - int(offset)
	default:
		return 0, fmt.Errorf("rbytes: invalid whence")
	}
	if r.c < 0 {
		return 0, fmt.Errorf("rbytes: negative position")
	}
	return int64(r.c), nil
}
