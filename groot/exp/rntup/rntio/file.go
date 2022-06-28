// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntio

import "io"

type PageFileReader struct {
	r io.ReaderAt
}

func (pfr *PageFileReader) ReadPage(p *Page) error {
	panic("not implemented")
}

var (
	_ PageReader = (*PageFileReader)(nil)
)
