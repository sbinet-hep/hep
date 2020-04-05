// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntup

import (
	"fmt"

	"go-hep.org/x/hep/groot/riofs"
)

type Reader struct {
	f  *riofs.File
	nt *NTuple
}

func Open(fname, nt string) (*Reader, error) {
	f, err := riofs.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("rntup: could not open ROOT file: %w", err)
	}
	defer func() {
		if err != nil {
			_ = f.Close()
		}
	}()

	obj, err := riofs.Dir(f).Get(nt)
	if err != nil {
		return nil, fmt.Errorf("rntup: could not get NTuple anchor: %w", err)
	}

	ntup, ok := obj.(*NTuple)
	if !ok {
		return nil, fmt.Errorf("rntup: object %q is not an NTuple (%T)", nt, obj)
	}

	return &Reader{
		f:  f,
		nt: ntup,
	}, nil
}

func (r *Reader) Close() error {
	return r.f.Close()
}
