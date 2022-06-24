// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntup

type Anchor struct {
	// The ROOT streamer info checksum. Older RNTuple versions used class version 0 and a serialized checksum,
	// now we use class version 3 and "promote" the checksum as a class member
	checksum int32

	// Allows for evolving the struct in future versions
	version uint32

	// Allows for skipping the struct
	size uint32 // sizeof(Anchor);

	// The file offset of the header excluding the TKey part
	seekHeader uint64

	// The size of the compressed ntuple header
	nbytesHeader uint32

	// The size of the uncompressed ntuple header
	lenHeader uint32

	// The file offset of the footer excluding the TKey part
	seekFooter uint64

	// The size of the compressed ntuple footer
	nbytesFooter uint32

	// The size of the uncompressed ntuple footer
	lenFooter uint32

	// Currently unused, reserved for later use
	reserved uint64
}
