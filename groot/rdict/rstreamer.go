// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"sync"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
)

type Unmarshaler interface {
	root.Object
	rbytes.Unmarshaler
}

type Marshaler interface {
	root.Object
	rbytes.Marshaler
}

// RStream reads a value from the underlying buffer.
func RStream(r *rbytes.RBuffer, o Unmarshaler) error {
	if r.Err() != nil {
		return r.Err()
	}

	cls := o.Class()
	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(cls)

	sikey := streamerKey{cls, vers}
	streamers.RLock()
	streamer, ok := streamers.db[sikey]
	streamers.RUnlock()
	if !ok {
		si, err := r.StreamerInfo(cls, int(vers))
		if err != nil {
			r.SetErr(errors.Wrapf(err, "rdict: could not find streamer info for %q v=%d", cls, vers))
			return r.Err()
		}
		streamer, err = buildStreamer(r, si)
		if err != nil {
			r.SetErr(errors.Wrapf(err, "rdict: could not build streamer for %q v=%d", cls, vers))
			return r.Err()
		}
		streamers.Lock()
		streamers.db[sikey] = streamer
		streamers.Unlock()
	}

	r.CheckByteCount(pos, bcnt, beg, cls)

	return r.Err()
}

type streamerKey struct {
	name string
	vers int16
}

type streamerVal struct {
	r rstreamerFunc
	w wstreamerFunc
}

var streamers = struct {
	sync.RWMutex
	db map[streamerKey]streamerVal
}{db: make(map[streamerKey]streamerVal)}

type rstreamerFunc func(r *rbytes.RBuffer, o Unmarshaler) error
type wstreamerFunc func(w *rbytes.WBuffer, o Marshaler) error

func buildStreamer(ctx rbytes.StreamerInfoContext, si rbytes.StreamerInfo) (streamerVal, error) {
	panic("not implemented")
}
