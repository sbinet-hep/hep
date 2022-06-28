// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnttypes // import "go-hep.org/x/hep/groot/exp/rntup/rnttypes"

type Kind uint16

const (
	Index64 Kind = 0x01 // Mother columns of (nested) collections, counting is relative to the cluster
	Index32 Kind = 0x02 // Mother columns of (nested) collections, counting is relative to the cluster
	Switch  Kind = 0x03 // Lower 44 bits like kIndex64, higher 20 bits are a dispatch tag to a column ID
	Byte    Kind = 0x04 // An uninterpreted byte, e.g. part of a blob
	Char    Kind = 0x05 // ASCII character
	Bit     Kind = 0x06 // Boolean value
	Real64  Kind = 0x07 // IEEE-754 double precision float
	Real32  Kind = 0x08 // IEEE-754 single precision float
	Real16  Kind = 0x09 // IEEE-754 half precision float
	Int64   Kind = 0x0A // Two's complement, little-endian 8 byte integer
	Int32   Kind = 0x0B // Two's complement, little-endian 4 byte integer
	Int16   Kind = 0x0C // Two's complement, little-endian 2 byte integer
	Int8    Kind = 0x0D // Two's complement, 1 byte integer

	SplitIndex64 Kind = 0x0E // Like Index64 but pages are stored in split + delta encoding
	SplitIndex32 Kind = 0x0F // Like Index32 but pages are stored in split + delta encoding
	SplitReal64  Kind = 0x10 // Like Real64 but pages are stored in split encoding
	SplitReal32  Kind = 0x11 // Like Real32 but pages are stored in split encoding
	SplitReal16  Kind = 0x12 // Like Real16 but pages are stored in split encoding
	SplitInt64   Kind = 0x13 // Like Int64 but pages are stored in split encoding
	SplitInt32   Kind = 0x14 // Like Int32 but pages are stored in split encoding
	SplitInt16   Kind = 0x15 // Like Int16 but pages are stored in split encoding
)

type ColumnFlag uint32

const (
	ColumnSortedIncr  ColumnFlag = 0x01 // Elements in the column are sorted (monotonically increasing)
	ColumnSortedDecr  ColumnFlag = 0x02 // Elements in the column are sorted (monotonically decreasing)
	ColumnNonNegative ColumnFlag = 0x04 // Elements have only non-negative values
)

type FieldFlag uint16

const (
	FlagRepetitive FieldFlag = 0x01 // Repetitive field, i.e. for every entry n copies of the field are stored
	FlagAlias      FieldFlag = 0x02 // Alias field, the columns referring to this field are alias columns
)

type FieldRole uint16

const (
	RoleLeaf       FieldRole = 0x00 // Leaf field in the schema tree
	RoleCollection FieldRole = 0x01 // The field is the mother of a collection (e.g., a vector)
	RoleRecord     FieldRole = 0x02 // The field is the mother of a record (e.g., a struct)
	RoleUnion      FieldRole = 0x03 // The field is the mother of a variant (e.g., a union)
	RoleReference  FieldRole = 0x04 // The field is a reference (pointer), TODO
)
