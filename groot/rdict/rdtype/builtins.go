// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdtype

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
)

type Bool struct {
	ptr *bool
}

func (Bool) Type() reflect.Type       { return boolType }
func (v *Bool) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Bool) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Bool) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadBool()
	return r.Err()
}

type Int8 struct {
	ptr *int8
}

func (Int8) Type() reflect.Type       { return int8Type }
func (v *Int8) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Int8) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Int8) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadI8()
	return r.Err()
}

type Int16 struct {
	ptr *int16
}

func (Int16) Type() reflect.Type       { return int16Type }
func (v *Int16) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Int16) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Int16) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadI16()
	return r.Err()
}

type Int32 struct {
	ptr *int32
}

func (Int32) Type() reflect.Type       { return int32Type }
func (v *Int32) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Int32) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Int32) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadI32()
	return r.Err()
}

type Int64 struct {
	ptr *int64
}

func (Int64) Type() reflect.Type       { return int64Type }
func (v *Int64) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Int64) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Int64) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadI64()
	return r.Err()
}

type Uint8 struct {
	ptr *uint8
}

func (Uint8) Type() reflect.Type       { return uint8Type }
func (v *Uint8) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Uint8) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Uint8) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadU8()
	return r.Err()
}

type Uint16 struct {
	ptr *uint16
}

func (Uint16) Type() reflect.Type       { return uint16Type }
func (v *Uint16) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Uint16) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Uint16) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadU16()
	return r.Err()
}

type Uint32 struct {
	ptr *uint32
}

func (Uint32) Type() reflect.Type       { return uint32Type }
func (v *Uint32) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Uint32) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Uint32) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadU32()
	return r.Err()
}

type Uint64 struct {
	ptr *uint64
}

func (Uint64) Type() reflect.Type       { return uint64Type }
func (v *Uint64) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Uint64) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Uint64) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadU64()
	return r.Err()
}

type Float16 struct {
	ptr *root.Float16
	elm rbytes.StreamerElement
}

func (Float16) Type() reflect.Type       { return float16Type }
func (v *Float16) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Float16) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Float16) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadF16(v.elm)
	return r.Err()
}

type Double32 struct {
	ptr *root.Double32
	elm rbytes.StreamerElement
}

func (Double32) Type() reflect.Type       { return double32Type }
func (v *Double32) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Double32) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Double32) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadD32(v.elm)
	return r.Err()
}

type Float32 struct {
	ptr *float32
}

func (Float32) Type() reflect.Type       { return float32Type }
func (v *Float32) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Float32) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Float32) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadF32()
	return r.Err()
}

type Float64 struct {
	ptr *float64
}

func (Float64) Type() reflect.Type       { return float64Type }
func (v *Float64) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Float64) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Float64) RStream(r *rbytes.RBuffer) error {
	*v.ptr = r.ReadF64()
	return r.Err()
}

type Complex64 struct {
	ptr *complex64
}

func (Complex64) Type() reflect.Type       { return complex64Type }
func (v *Complex64) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Complex64) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Complex64) RStream(r *rbytes.RBuffer) error {
	re := r.ReadF32()
	im := r.ReadF32()
	*v.ptr = complex(re, im)
	return r.Err()
}

type Complex128 struct {
	ptr *complex128
}

func (Complex128) Type() reflect.Type       { return complex128Type }
func (v *Complex128) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }

func (v *Complex128) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *Complex128) RStream(r *rbytes.RBuffer) error {
	re := r.ReadF64()
	im := r.ReadF64()
	*v.ptr = complex(re, im)
	return r.Err()
}

type String struct {
	ptr *string
}

func (String) Type() reflect.Type       { return stringType }
func (v *String) rvalue() reflect.Value { return reflect.ValueOf(v.ptr) }
func (v String) Len() int               { return len(*v.ptr) }

func (v *String) UnmarshalROOT(r *rbytes.RBuffer) error {
	return v.RStream(r)
}

func (v *String) RStream(r *rbytes.RBuffer) error {
	// FIXME(sbinet): std::string/TString/char* ambiguity
	*v.ptr = r.ReadString()
	return r.Err()
}

var (
	_ Value = (*Bool)(nil)

	_ Value = (*Int8)(nil)
	_ Value = (*Int16)(nil)
	_ Value = (*Int32)(nil)
	_ Value = (*Int64)(nil)

	_ Value = (*Uint8)(nil)
	_ Value = (*Uint16)(nil)
	_ Value = (*Uint32)(nil)
	_ Value = (*Uint64)(nil)

	_ Value = (*Float16)(nil)
	_ Value = (*Double32)(nil)
	_ Value = (*Float32)(nil)
	_ Value = (*Float64)(nil)

	_ Value = (*Complex64)(nil)
	_ Value = (*Complex128)(nil)

	_ Value = (*String)(nil)
)
