// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntup_test

import (
	"testing"

	"go-hep.org/x/hep/groot/exp/rntup"
)

func TestReader(t *testing.T) {
	r, err := rntup.Open("../../testdata/ntpl001_staff.root", "Staff")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer r.Close()
}
