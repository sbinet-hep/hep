// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rarrow // import "go-hep.org/x/hep/groot/rarrow"

import (
	"fmt"

	"github.com/apache/arrow/go/arrow"
	"github.com/pkg/errors"

	"go-hep.org/x/hep/groot/internal/rmeta"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
)

func SchemaFrom(sictx rbytes.StreamerInfoContext, si rbytes.StreamerInfo) *arrow.Schema {
	var (
		bld    = builder{si: si, sictx: sictx}
		md     = bld.genMetadata(si)
		fields = make([]arrow.Field, 0, len(si.Elements()))
	)
	fmt.Printf("gen-schema: %s %q\n", si.Name(), si.Title())
	for _, elem := range si.Elements() {
		fmt.Printf("--> %s (%T)\n", elem.Name(), elem)
	}
	for i, elem := range si.Elements() {
		fields = append(fields, bld.genField(elem, i))
	}
	return arrow.NewSchema(fields, &md)
}

type builder struct {
	si    rbytes.StreamerInfo
	sictx rbytes.StreamerInfoContext
}

func (builder) genMetadata(si rbytes.StreamerInfo) arrow.Metadata {
	kv := map[string]string{
		"name":          si.Name(),
		"title":         si.Title(),
		"class":         si.Class(),
		"class_version": fmt.Sprintf("%d", si.ClassVersion()),
		"checksum":      fmt.Sprintf("%d", si.CheckSum()),
	}
	return arrow.MetadataFrom(kv)
}

func (b *builder) genField(se rbytes.StreamerElement, i int) arrow.Field {
	return arrow.Field{
		Name: se.Name(),
		Type: b.r2arr(se, i),
	}
}

type todoDataType struct {
	name string
	tid  arrow.Type
}

func (dt todoDataType) Name() string   { return dt.name }
func (dt todoDataType) ID() arrow.Type { return arrow.NULL }

var (
	_ arrow.DataType = (*todoDataType)(nil)
)

func (b builder) r2arr(se rbytes.StreamerElement, i int) arrow.DataType {
	switch se := se.(type) {
	case *rdict.StreamerBase:
		return arrow.StructOf(
			//			arrow.Field{Name: fmt.Sprintf("ROOT_base_%03d", i), Type: r2arr(sb, i)},
			arrow.Field{Name: fmt.Sprintf("ROOT_base_%03d_%s", i, se.Name()), Type: tobject},
		)

	case *rdict.StreamerBasicType:
		switch rt := se.Type(); rt {
		case rmeta.Counter:
			switch se.Size() {
			case 4:
				return arrow.PrimitiveTypes.Int32
			case 8:
				return arrow.PrimitiveTypes.Int64
			default:
				panic(errors.Errorf("rarrow: invalid rmeta.Counter size %d", se.Size()))
			}

		case rmeta.Char:
			return arrow.PrimitiveTypes.Int8
		case rmeta.Short:
			return arrow.PrimitiveTypes.Int16
		case rmeta.Int:
			return arrow.PrimitiveTypes.Int32
		case rmeta.Long, rmeta.Long64:
			return arrow.PrimitiveTypes.Int64

		case rmeta.CharStar, rmeta.LegacyChar:
			return arrow.BinaryTypes.Binary

		case rmeta.Float:
			return arrow.PrimitiveTypes.Float32
		case rmeta.Double:
			return arrow.PrimitiveTypes.Float64
		case rmeta.Float16:
			panic("float16 not supported") // FIXME(sbinet)
			return arrow.PrimitiveTypes.Float32
		case rmeta.Double32:
			panic("double32 not supported") // FIXME(sbinet)
			return arrow.PrimitiveTypes.Float64

		case rmeta.UChar:
			return arrow.PrimitiveTypes.Uint8
		case rmeta.UShort:
			return arrow.PrimitiveTypes.Uint16
		case rmeta.UInt:
			return arrow.PrimitiveTypes.Uint32
		case rmeta.ULong, rmeta.ULong64:
			return arrow.PrimitiveTypes.Uint64

		case rmeta.Bool:
			return arrow.FixedWidthTypes.Boolean

		case rmeta.OffsetL + rmeta.Char:
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Int8)
		case rmeta.OffsetL + rmeta.Short:
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Int16)
		case rmeta.OffsetL + rmeta.Int:
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Int32)
		case rmeta.OffsetL + rmeta.Long, rmeta.OffsetL + rmeta.Long64:
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Int64)
		case rmeta.OffsetL + rmeta.UChar:
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Uint8)
		case rmeta.OffsetL + rmeta.UShort:
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Uint16)
		case rmeta.OffsetL + rmeta.UInt:
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Uint32)
		case rmeta.OffsetL + rmeta.ULong, rmeta.OffsetL + rmeta.ULong64:
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Uint64)

		case rmeta.OffsetL + rmeta.Float:
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Float32)
		case rmeta.OffsetL + rmeta.Double:
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Float64)
		case rmeta.OffsetL + rmeta.Float16:
			panic("float16 not supported") // FIXME(sbinet)
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Float32)
		case rmeta.OffsetL + rmeta.Double32:
			panic("double32 not supported") // FIXME(sbinet)
			return arrayOf(se.ArrayLen(), arrow.PrimitiveTypes.Float64)

		case rmeta.OffsetL + rmeta.Bool:
			return arrayOf(se.ArrayLen(), arrow.FixedWidthTypes.Boolean)

		default:
			panic(errors.Errorf("rarrow: invalid StreamerBasicType: %#v", se))
		}

	case *rdict.StreamerBasicPointer:
		switch se.Type() {
		case rmeta.OffsetP + rmeta.Char:
			return sliceOf(arrow.PrimitiveTypes.Int8)
		case rmeta.OffsetP + rmeta.Short:
			return sliceOf(arrow.PrimitiveTypes.Int16)
		case rmeta.OffsetP + rmeta.Int:
			return sliceOf(arrow.PrimitiveTypes.Int32)
		case rmeta.OffsetP + rmeta.Long, rmeta.OffsetP + rmeta.Long64:
			return sliceOf(arrow.PrimitiveTypes.Int64)
		case rmeta.OffsetP + rmeta.Float:
			return sliceOf(arrow.PrimitiveTypes.Float32)
		case rmeta.OffsetP + rmeta.Float16:
			panic("float16 not supported") // FIXME(sbinet)
			return sliceOf(arrow.PrimitiveTypes.Float32)
		case rmeta.OffsetP + rmeta.Double32:
			panic("double32 not supported") // FIXME(sbinet)
			return sliceOf(arrow.PrimitiveTypes.Float64)
		case rmeta.OffsetP + rmeta.Double:
			return sliceOf(arrow.PrimitiveTypes.Float64)
		case rmeta.OffsetP + rmeta.UChar, rmeta.OffsetP + rmeta.CharStar:
			return sliceOf(arrow.PrimitiveTypes.Uint8)
		case rmeta.OffsetP + rmeta.UShort:
			return sliceOf(arrow.PrimitiveTypes.Uint16)
		case rmeta.OffsetP + rmeta.UInt, rmeta.OffsetP + rmeta.Bits:
			return sliceOf(arrow.PrimitiveTypes.Uint32)
		case rmeta.OffsetP + rmeta.ULong, rmeta.OffsetP + rmeta.ULong64:
			return sliceOf(arrow.PrimitiveTypes.Uint64)
		case rmeta.OffsetP + rmeta.Bool:
			return sliceOf(arrow.FixedWidthTypes.Boolean)
		default:
			panic(errors.Errorf("rarrow: invalid StreamerBasicPointer: %#v", se))
		}
	case *rdict.StreamerLoop:
		panic(errors.Errorf("rarrow: StreamerLoop not supported %#v", se))

	case *rdict.StreamerObject:
		sio, err := b.sictx.StreamerInfo(se.TypeName())
		if err != nil {
			panic(errors.Wrapf(err, "could not find StreamerInfo for StreamerObject"))
		}
		return b.genType(sio)
	case *rdict.StreamerObjectPointer:

	case *rdict.StreamerObjectAny:

	case *rdict.StreamerString:
		return arrow.BinaryTypes.String
	case *rdict.StreamerSTL:
	case *rdict.StreamerSTLstring:
		return arrow.BinaryTypes.String
	case *rdict.StreamerArtificial:
	default:

	}

	switch rt := se.Type(); rt {
	case rmeta.Base:
		sb := se.(*rdict.StreamerBase)
		return arrow.StructOf(
			//			arrow.Field{Name: fmt.Sprintf("ROOT_base_%03d", i), Type: r2arr(sb, i)},
			arrow.Field{Name: fmt.Sprintf("ROOT_base_%03d_%s", i, sb.Name()), Type: tobject},
		)
	case rmeta.Char:
		return arrow.PrimitiveTypes.Int8
	case rmeta.Short:
		return arrow.PrimitiveTypes.Int16
	case rmeta.Int:
		return arrow.PrimitiveTypes.Int32
	case rmeta.Long, rmeta.Long64:
		return arrow.PrimitiveTypes.Int64
	case rmeta.Counter:
		return arrow.PrimitiveTypes.Int32
	case rmeta.CharStar, rmeta.LegacyChar:
		return arrow.BinaryTypes.Binary

	case rmeta.Float:
		return arrow.PrimitiveTypes.Float32
	case rmeta.Double:
		return arrow.PrimitiveTypes.Float64
	case rmeta.Double32:
		panic("double32 not supported") // FIXME(sbinet)
		return arrow.PrimitiveTypes.Float64

	case rmeta.UChar:
		return arrow.PrimitiveTypes.Uint8
	case rmeta.UShort:
		return arrow.PrimitiveTypes.Uint16
	case rmeta.UInt:
		return arrow.PrimitiveTypes.Uint32
	case rmeta.ULong, rmeta.ULong64:
		return arrow.PrimitiveTypes.Uint64

	case rmeta.TNamed:
		return tnamed
	default:
		panic(errors.Errorf("groot/rarrow: invalid ROOT type %d %#v", rt, se))
	}
}

func (b *builder) genType(si rbytes.StreamerInfo) arrow.DataType {
	return nil
}

func arrayOf(n int, dt arrow.DataType) arrow.DataType {
	// FIXME(sbinet)
	return arrow.ListOf(dt)
}

func sliceOf(dt arrow.DataType) arrow.DataType {
	return arrow.ListOf(dt)
}

var (
	tnamed  arrow.DataType
	tobject arrow.DataType
)

func init() {
	tobject = arrow.StructOf(
		arrow.Field{Name: "ID", Type: arrow.PrimitiveTypes.Uint32},
		arrow.Field{Name: "Bits", Type: arrow.PrimitiveTypes.Uint32},
	)

	tnamed = arrow.StructOf(
		arrow.Field{Name: "obj", Type: tobject},
		arrow.Field{Name: "name", Type: arrow.BinaryTypes.String},
		arrow.Field{Name: "title", Type: arrow.BinaryTypes.String},
	)

}
