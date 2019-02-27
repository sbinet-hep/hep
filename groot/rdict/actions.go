// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import "go-hep.org/x/hep/groot/rbytes"

type compInfo struct {
	offset   []int
	length   int
	elem     rbytes.StreamerElement
	class    string
	sinfo    rbytes.StreamerInfo
	streamer memberStreamer
}

type config struct {
	sinfo  rbytes.StreamerInfo // StreamerInfo for which the action is derived
	elem   rbytes.StreamerElement
	comp   *compInfo // compiled information
	offset []int     // offset within object
	length int       // number of elements in a fixed length array
}

// loopConfig is the interface for member-wise looping routines configuration.
type loopConfig interface {
}

type action struct {
	cfg  *config
	ract func(r *rbytes.RBuffer, recv interface{}, cfg *config)
	wact func(w *rbytes.RBuffer, recv interface{}, cfg *config)

	// TODO: vec-ptr-loop action (r/w)
	// TODO: loop-action (r/w)
}

type actions struct {
	sinfo   rbytes.StreamerInfo // StreamerInfo used to derive these actions
	loop    loopConfig          // if this is a bundle of memberwise streaming actions, this configures the looping
	actions []action
}

type memberStreamer struct{}
