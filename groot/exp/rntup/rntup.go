// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rntup contains types to handle RNTuple-related data.
package rntup // import "go-hep.org/x/hep/groot/exp/rntup"

import (
	"go-hep.org/x/hep/groot/exp/rntup/rnttypes"
)

type span struct {
	seek   uint64 // file offset of the span, excluding TKey part
	nbytes uint32 // size of compressed span
	length uint32 // size of uncompressed span
}

type header struct {
	vers uint16
	minv uint16

	flags   []uint64
	release uint32
	name    string
	descr   string
	library string

	fields  []fieldDescr
	cols    []colDescr
	aliases []colAlias
	extra   []colExtra

	crc32 uint32
}

type footer struct {
	vers uint16
	minv uint16

	flags []uint64
	hdr   uint32

	xhdrs     []extHeader
	colGroups []colGroup
	clInfos   []clusterInfo
	clGroups  []clusterGroup
	mdBlocks  []metaDataBlock

	crc32 uint32
}

type pageList struct {
	vers uint16
	minv uint16

	clusters []clusterDescr

	crc32 uint32
}

type frame struct {
	size uint32
	n    uint32
}

type fieldDescr struct {
	vers uint32 // field version
	typv uint32 // type version
	pfid uint32 // parent field ID
	role rnttypes.FieldRole
	flag rnttypes.FieldFlag
	nrep uint64

	fname string // field name
	tname string // type name
	alias string // type alias
	descr string // field description
}

type colDescr struct {
	kind    rnttypes.Kind
	bits    uint16
	fieldID uint32
	flags   rnttypes.ColumnFlag
}

type colAlias struct {
	physID  uint32 // physical column ID
	fieldID uint32 // field that needs to have the "alias field" flag set
}

type colExtra struct {
	typeFrom  uint32
	typeTo    uint32
	contentID uint32
	typeName  string
}

type extHeader struct {
	// TODO
}

type colGroup struct {
	// TODO
}

type clusterInfo struct {
	firstEntry uint64
	nentries   uint64
	colGrpID   int32 // -1 for "all columns"
}

type clusterGroup struct {
	n     uint32
	pages envelopeLink
}

type metaDataBlock struct {
}

type envelopeLink struct {
	size uint32 // unzipped size
	loc  locator
}

type locator struct {
	pos     uint64
	storage uint32
	url     string
}

type clusterDescr struct {
	columns []columnDescr
}

type columnDescr struct {
	pages  []pageDescr
	offset uint64
	compr  uint32
}

type pageDescr struct {
	nelem uint32
	loc   locator
}
