// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntup

// Page is a slice of a column that is mapped into memory.
type Page struct {
	col  int64 // Column id
	buf  []byte
	esz  int64 // element size in bytes
	nelm int64 // number of elements
	cap  int64 // capacity of the page in number of elements

	rng     uint64
	cluster ClusterInfo
}

func (p *Page) Contains(idx uint64) bool {
	return p.rng <= idx && idx < p.rng+uint64(p.nelm)
}

type ClusterInfo struct {
	id   uint64 // cluster number.
	ioff uint64 // first element index of the column in this cluster.
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
