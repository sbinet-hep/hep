// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rarrow_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rarrow"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/rdict"
	_ "go-hep.org/x/hep/groot/ztypes"
)

func TestSchema(t *testing.T) {
	for _, si := range sinfos {
		fmt.Printf("=================\n")
		fmt.Printf("schema for %q\n", si.Name())
		schema := rarrow.SchemaFrom(sictx, si)
		md := schema.Metadata()
		for i, k := range md.Keys() {
			fmt.Printf("%s: %q\n", k, md.Values()[i])
		}
		fmt.Printf("fields: %d\n", len(schema.Fields()))
		for _, f := range schema.Fields() {
			fmt.Printf("%s: %v\n", f.Name, f.Type.Name())
		}
	}
}

var (
	sictx  = &context{db: make(map[string]rbytes.StreamerInfo)}
	sinfos []rbytes.StreamerInfo
)

type context struct {
	db map[string]rbytes.StreamerInfo
}

func (ctx context) StreamerInfo(name string) (rbytes.StreamerInfo, error) {
	si, ok := ctx.db[name]
	if !ok {
		return nil, errors.Errorf("no such streamer info %q", name)
	}
	return si, nil
}

func init() {
	for _, si := range []rbytes.StreamerInfo{
		rdict.StreamerOf(sictx, reflect.TypeOf((*struct1)(nil)).Elem()),
		//		rdict.StreamerOf(sictx, reflect.TypeOf((*struct2)(nil)).Elem()),
		//		rdict.StreamerOf(sictx, reflect.TypeOf((*struct3)(nil)).Elem()),
		//		rdict.StreamerOf(sictx, reflect.TypeOf((*struct4)(nil)).Elem()),
	} {
		sinfos = append(sinfos, si)
	}
	if true {
		return
	}

	data, err := ioutil.ReadFile("../testdata/tlist-tsi.dat")
	if err != nil {
		log.Fatal(err)
	}

	r := rbytes.NewRBuffer(data, nil, 0, nil)
	var lst rcont.List
	err = lst.UnmarshalROOT(r)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < lst.Len(); i++ {
		switch obj := lst.At(i).(type) {
		case *rdict.StreamerInfo:
			sinfos = append(sinfos, obj)
			sictx.db[obj.Name()] = obj
		}
	}
}

// 	rdict.NewStreamerInfo("MyClass", []rbytes.StreamerElement{
// 		rdict.NewStreamerSTL("DataVector<int>", rmeta.STLvector, rmeta.Int),
// 	}),
// }

var (
	_ rbytes.StreamerInfoContext = (*context)(nil)
)

type struct1 struct {
	Name    string
	Bool    bool
	I8      int8
	I16     int16
	I32     int32
	I64     int64
	U8      uint8
	U16     uint16
	U32     uint32
	U64     uint64
	F32     float32
	F64     float64
	Float64 float64 `groot:"Cxx::MyFloat64"`
}

type struct2 struct {
	V1 struct1
}

type struct3 struct {
	Names [10]string
	Bools [10]bool
	I8s   [10]int8
	I16s  [10]int16
	I32s  [10]int32
	I64s  [10]int64
	U8s   [10]uint8
	U16s  [10]uint16
	U32s  [10]uint32
	U64s  [10]uint64
	F32s  [10]float32
	F64s  [10]float64
	S1s   [10]struct1
}

type struct4 struct {
	Names []string
	Bools []bool
	I8s   []int8
	I16s  []int16
	I32s  []int32
	I64s  []int64
	U8s   []uint8
	U16s  []uint16
	U32s  []uint32
	U64s  []uint64
	F32s  []float32
	F64s  []float64
	S1s   []struct1
}

type struct5 struct {
	I32 *int32
	F64 *float64
	S1  *struct1
}
