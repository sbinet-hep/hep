// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntup

// Page is a slice of a column that is mapped into memory.
type Page struct {
	col  int // Column id
	buf  []byte
	cap  int
	esz  int // element size in bytes
	nelm int // number of elements

	rng     int
	cluster ClusterInfo
}

func (p *Page) Contains(idx int) bool {
	return idx >= p.rng && idx < p.rng+p.nelm
}

type ClusterInfo struct {
	id   int // cluster number.
	ioff int // first element index of the column in this cluster.
}

type ColumnHandle struct {
	id  int
	col *Column
}

type Column struct{}

type PageReader interface {
	Attach()
	Clone() PageReader
	Read(*Page) error
}
