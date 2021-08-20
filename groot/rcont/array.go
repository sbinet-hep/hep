// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcont

import (
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rvers"
)

type Array[T any] struct {
	Data []T
}

func (*Array[T]) RVersion() int16 {
	return rvers.ArrayC
}

// Class returns the ROOT class name.
func (*Array[T]) Class() string {
	return "TArray[T any]"
}

func (arr *Array[T]) Len() int {
	return len(arr.Data)
}

func (arr *Array[T]) At(i int) T {
	return arr.Data[i]
}

func (arr *Array[T]) Get(i int) interface{} {
	return arr.Data[i]
}

func (arr *Array[T]) Set(i int, v interface{}) {
	arr.Data[i] = v.(T)
}

func (arr *Array[T]) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	rbytes.WriteArray(w, arr.Data)

	return int(w.Pos() - pos), w.Err()
}

func (arr *Array[T]) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	n := int(r.ReadI32())
	arr.Data = rbytes.Resize[T](arr.Data, n)
	rbytes.ReadArray(r, arr.Data)

	return r.Err()
}
