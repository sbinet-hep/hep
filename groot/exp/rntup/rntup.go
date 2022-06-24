// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rntup contains types to handle RNTuple-related data.
package rntup // import "go-hep.org/x/hep/groot/exp/rntup"

type fieldFlag uint16

const (
	fieldFlagRepetitive fieldFlag = 0x01 // Repetitive field, i.e. for every entry n copies of the field are stored
	fieldFlagAlias      fieldFlag = 0x02 // Alias field, the columns referring to this field are alias columns
)

type fieldRole uint16

const (
	fieldRoleLeaf       fieldRole = 0x00 // Leaf field in the schema tree
	fieldRoleCollection fieldRole = 0x01 // The field is the mother of a collection (e.g., a vector)
	fieldRoleRecord     fieldRole = 0x02 // The field is the mother of a record (e.g., a struct)
	fieldRoleUnion      fieldRole = 0x03 // The field is the mother of a variant (e.g., a union)
	fieldRoleReference  fieldRole = 0x04 // The field is a reference (pointer), TODO
)

type span struct {
	seek   uint64 // file offset of the span, excluding TKey part
	nbytes uint32 // size of compressed span
	length uint32 // size of uncompressed span
}

// type envelope struct {
// 	vers  uint16
// 	minv  uint16
// 	crc32 uint32
// }

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
	role fieldRole
	flag fieldFlag
	nrep uint64

	fname string // field name
	tname string // type name
	alias string // type alias
	descr string // field description
}

type colKind uint16

const (
	colIndex64 colKind = 0x01 // Mother columns of (nested) collections, counting is relative to the cluster
	colIndex32 colKind = 0x02 // Mother columns of (nested) collections, counting is relative to the cluster
	colSwitch  colKind = 0x03 // Lower 44 bits like kIndex64, higher 20 bits are a dispatch tag to a column ID
	colByte    colKind = 0x04 // An uninterpreted byte, e.g. part of a blob
	colChar    colKind = 0x05 // ASCII character
	colBit     colKind = 0x06 // Boolean value
	colReal64  colKind = 0x07 // IEEE-754 double precision float
	colReal32  colKind = 0x08 // IEEE-754 single precision float
	colReal16  colKind = 0x09 // IEEE-754 half precision float
	colInt64   colKind = 0x0A // Two's complement, little-endian 8 byte integer
	colInt32   colKind = 0x0B // Two's complement, little-endian 4 byte integer
	colInt16   colKind = 0x0C // Two's complement, little-endian 2 byte integer
	colInt8    colKind = 0x0D // Two's complement, 1 byte integer

	colSplitIndex64 colKind = 0x0E // Like Index64 but pages are stored in split + delta encoding
	colSplitIndex32 colKind = 0x0F // Like Index32 but pages are stored in split + delta encoding
	colSplitReal64  colKind = 0x10 // Like Real64 but pages are stored in split encoding
	colSplitReal32  colKind = 0x11 // Like Real32 but pages are stored in split encoding
	colSplitReal16  colKind = 0x12 // Like Real16 but pages are stored in split encoding
	colSplitInt64   colKind = 0x13 // Like Int64 but pages are stored in split encoding
	colSplitInt32   colKind = 0x14 // Like Int32 but pages are stored in split encoding
	colSplitInt16   colKind = 0x15 // Like Int16 but pages are stored in split encoding
)

type colFlag uint32

const (
	colSortedIncr  colFlag = 0x01 // Elements in the column are sorted (monotonically increasing)
	colSortedDecr  colFlag = 0x02 // Elements in the column are sorted (monotonically decreasing)
	colNonNegative colFlag = 0x04 // Elements have only non-negative values
)

type colDescr struct {
	kind    colKind
	bits    uint16
	fieldID uint32
	flags   colFlag
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
