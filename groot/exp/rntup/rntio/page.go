// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntio

type PageReader interface {
	ReadPage(p *Page) error
}

type PageWriter interface {
	WritePage(p *Page) error
}

type PageAllocator interface {
	PageAlloc() *Page
	PageRealloc(p *Page) *Page
	PageFree(p *Page)
}

type Page struct {
	colID ColumnID
	buf   []byte
	elmt  struct {
		size uint32 // size of an element
		n    uint32 // number of elements in this page
		cap  uint32 // capacity of the page in number of elements
	}
	first   uint64 // first entry number in this page
	cluster ClusterInfo
}

type DescrID uint64
type ColumnID int64
type NtupleSize uint64

type ClusterInfo struct {
	id        DescrID    // cluster number
	idxOffset NtupleSize // first element index of the column in this cluster
}
