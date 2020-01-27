// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdtype

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"golang.org/x/xerrors"
)

type Value interface {
	Type() reflect.Type

	rvalue() reflect.Value

	rbytes.Unmarshaler
	// rbytes.Marshaler

	RStream(r *rbytes.RBuffer) error
	// WStream(w *rbytes.WBuffer) (int, error)
}

var (
	boolType   = reflect.TypeOf(false)
	uint8Type  = reflect.TypeOf(uint8(0))
	uint16Type = reflect.TypeOf(uint16(0))
	uint32Type = reflect.TypeOf(uint32(0))
	uint64Type = reflect.TypeOf(uint64(0))
	int8Type   = reflect.TypeOf(int8(0))
	int16Type  = reflect.TypeOf(int16(0))
	int32Type  = reflect.TypeOf(int32(0))
	int64Type  = reflect.TypeOf(int64(0))

	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))

	complex64Type  = reflect.TypeOf(complex(float32(0), float32(0)))
	complex128Type = reflect.TypeOf(complex(float64(0), float64(0)))

	stringType = reflect.TypeOf("")

	float16Type  = reflect.TypeOf(root.Float16(0))
	double32Type = reflect.TypeOf(root.Double32(0))
)

func New(typ reflect.Type, sictx rbytes.StreamerInfoContext) Value {
	v := reflect.New(typ)
	return ValueOf(v.Interface(), sictx)
}

func ValueOf(ptr interface{}, sictx rbytes.StreamerInfoContext) Value {
	typ := reflect.TypeOf(ptr).Elem()
	switch typ {
	case boolType:
		return &Bool{ptr.(*bool)}

	case uint8Type:
		return &Uint8{ptr.(*uint8)}

	case uint16Type:
		return &Uint16{ptr.(*uint16)}

	case uint32Type:
		return &Uint32{ptr.(*uint32)}

	case uint64Type:
		return &Uint64{ptr.(*uint64)}

	case int8Type:
		return &Int8{ptr.(*int8)}

	case int16Type:
		return &Int16{ptr.(*int16)}

	case int32Type:
		return &Int32{ptr.(*int32)}

	case int64Type:
		return &Int64{ptr.(*int64)}

	case float16Type:
		// FIXME(sbinet): how to receive rbytes.StreamerElement ?
		return &Float16{ptr: ptr.(*root.Float16)}

	case double32Type:
		// FIXME(sbinet): how to receive rbytes.StreamerElement ?
		return &Double32{ptr: ptr.(*root.Double32)}

	case float32Type:
		return &Float32{ptr.(*float32)}

	case float64Type:
		return &Float64{ptr.(*float64)}

	case complex64Type:
		return &Complex64{ptr.(*complex64)}

	case complex128Type:
		return &Complex128{ptr.(*complex128)}

	case stringType:
		return &String{ptr.(*string)}
	}

	switch typ.Kind() {
	default:
		panic(xerrors.Errorf("invalid type %T", ptr))

	case reflect.Array:
		return &Array{reflect.ValueOf(ptr)}

	case reflect.Slice:
		return &Slice{reflect.ValueOf(ptr)}

	case reflect.Struct:
		// FIXME(sbinet): fetch class name + version
		return &Struct{ptr: reflect.ValueOf(ptr), cls: typ.Name(), ver: -1}

	case reflect.Map:
		// FIXME(sbinet): fetch class name + version
		return &Map{ptr: reflect.ValueOf(ptr), cls: typ.Name(), ver: -1}
	}
}

type Array struct {
	ptr reflect.Value
}

func (v Array) Type() reflect.Type     { return v.ptr.Elem().Type() }
func (v *Array) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }
func (v Array) Len() int               { return v.ptr.Elem().Len() }
func (v Array) Index(i int) Value      { return ValueOf(v.ptr.Elem().Index(i).Addr().Interface(), nil) }

func (v *Array) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Array) RStream(r *rbytes.RBuffer) error {
	for i := 0; i < v.Len(); i++ {
		_ = v.Index(i).RStream(r)
	}
	return r.Err()
}

type Slice struct {
	ptr reflect.Value
}

func (v Slice) Type() reflect.Type     { return v.ptr.Elem().Type() }
func (v *Slice) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }
func (v Slice) Len() int               { return v.ptr.Elem().Len() }
func (v Slice) Index(i int) Value      { return ValueOf(v.ptr.Elem().Index(i).Addr().Interface(), nil) }

func (v *Slice) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Slice) RStream(r *rbytes.RBuffer) error {
	for i := 0; i < v.Len(); i++ {
		_ = v.Index(i).RStream(r)
	}
	return r.Err()
}

type Struct struct {
	ptr reflect.Value
	cls string
	ver int16
}

func (v Struct) Type() reflect.Type     { return v.ptr.Elem().Type() }
func (v *Struct) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }
func (v Struct) Field(i int) Value      { return ValueOf(v.ptr.Elem().Field(i).Addr().Interface(), nil) }

func (v *Struct) UnmarshalROOT(r *rbytes.RBuffer) error {
	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(v.cls)
	if vers != v.ver {
		r.SetErr(xerrors.Errorf("rdtype: inconsistent ROOT version for %q (got=%d, want=%d)", v.cls, vers, v.ver))
		return r.Err()
	}

	err := v.RStream(r)
	if err != nil {
		return xerrors.Errorf("rdtype: could not stream-out %q: %w", v.cls, err)
	}

	r.CheckByteCount(pos, bcnt, beg, v.cls)
	return r.Err()
}

func (v *Struct) RStream(r *rbytes.RBuffer) error {
	rt := v.Type()
	nfields := rt.NumField()
	for i := 0; i < nfields; i++ {
		_ = v.Field(i).RStream(r)
	}
	return r.Err()
}

type Map struct {
	ptr reflect.Value
	cls string
	ver int16
}

func (v Map) Type() reflect.Type     { return v.ptr.Elem().Type() }
func (v *Map) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Map) UnmarshalROOT(r *rbytes.RBuffer) error {
	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(v.cls)
	if vers != v.ver {
		r.SetErr(xerrors.Errorf("rdtype: inconsistent ROOT version for %q (got=%d, want=%d)", v.cls, vers, v.ver))
		return r.Err()
	}

	err := v.RStream(r)
	if err != nil {
		return xerrors.Errorf("rdtype: could not stream-out %q: %w", v.cls, err)
	}

	r.CheckByteCount(pos, bcnt, beg, v.cls)
	return r.Err()
}

func (v *Map) RStream(r *rbytes.RBuffer) error {
	rt := v.Type()
	n := int(r.ReadI32())

	kt := rt.Key()
	keys := New(kt, nil).(*Slice)
	keys.ptr.Elem().Set(reflect.MakeSlice(kt, n, n))

	err := keys.RStream(r)
	if err != nil {
		return xerrors.Errorf("rdtype: could not stream-out keys for %q: %w", v.cls, err)
	}

	// FIXME(sbinet): do not create temporary slices?
	vt := rt.Elem()
	vals := New(vt, nil).(*Slice)
	vals.ptr.Elem().Set(reflect.MakeSlice(vt, n, n))

	err = vals.RStream(r)
	if err != nil {
		return xerrors.Errorf("rdtype: could not stream-out vals for %q: %w", v.cls, err)
	}

	for i := 0; i < n; i++ {
		v.ptr.SetMapIndex(keys.ptr.Elem().Index(i), vals.ptr.Elem().Index(i))
	}

	return r.Err()
}

var (
	_ Value = (*Array)(nil)
	_ Value = (*Slice)(nil)
	_ Value = (*Struct)(nil)
	_ Value = (*Map)(nil)
)
