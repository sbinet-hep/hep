// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntmodel // import "go-hep.org/x/hep/groot/exp/rntup/rntmodel"

type Schema struct {
	id uint64 // unique schema ID

	root rootField
}

type Field interface {
	isField()
}

type rootField struct{}

func (rootField) isField() {}

var (
	_ Field = (*rootField)(nil)
)
