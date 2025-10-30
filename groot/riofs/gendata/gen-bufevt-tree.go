// Copyright ©2025 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"flag"
	"log"

	"go-hep.org/x/hep/groot/internal/rtests"
)

var (
	root  = flag.String("f", "test-bufevt.root", "output ROOT file")
	split = flag.Int("split", 0, "default split-level for TTree")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentree", []byte(script), *root, *split)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
#include <string.h>
#include <stdio.h>

const int ARRAYSZ  = 10;
const int MAXSLICE = 20;
const int MAXSTR   = 32;

#define OFFSET 0

struct __attribute__((packed)) Event {
	bool     Bool;
	char     Str[6];
	int8_t   I8;
	int16_t  I16;
	int32_t  I32;
	int64_t  I64;
	int64_t  G64;
	uint8_t  U8;
	uint16_t U16;
	uint32_t U32;
	uint64_t U64;
	uint64_t UGG;
	float    F32;
	double   F64;

	bool     ArrayBs[ARRAYSZ];
	int8_t   ArrayI8[ARRAYSZ];
	int16_t  ArrayI16[ARRAYSZ];
	int32_t  ArrayI32[ARRAYSZ];
	int64_t  ArrayI64[ARRAYSZ];
	int64_t  ArrayG64[ARRAYSZ];
	uint8_t  ArrayU8[ARRAYSZ];
	uint16_t ArrayU16[ARRAYSZ];
	uint32_t ArrayU32[ARRAYSZ];
	uint64_t ArrayU64[ARRAYSZ];
	uint64_t ArrayUGG[ARRAYSZ];
	float    ArrayF32[ARRAYSZ];
	double   ArrayF64[ARRAYSZ];

};

void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;
	int evtmax = 10;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "my tree title");

	Event e;
	t->Branch("Event",&e,
		"B/O"
		":Str[6]/C"
		":I8/B:I16/S:I32/I:I64/L:G64/G"
		":U8/b:U16/s:U32/i:U64/l:UGG/g"
		":F32/F:F64/D"

		// static arrays
		":ArrBs[10]/O"
		":ArrI8[10]/B:ArrI16[10]/S:ArrI32[10]/I:ArrI64[10]/L:ArrG64[10]/G"
		":ArrU8[10]/b:ArrU16[10]/s:ArrU32[10]/i:ArrU64[10]/l:ArrUGG[10]/g"
		":ArrF32[10]/F:ArrF64[10]/D"
	);

	for (int j = 0; j != evtmax; j++) {
		int i = j + OFFSET;
		e.Bool = (i % 2) == 0;
		strncpy(e.Str, TString::Format("str-%d\0", i).Data(), 32);
		e.I8  = -i;
		e.I16 = -i;
		e.I32 = -i;
		e.I64 = -i;
		e.G64 = -i;
		e.U8  = i;
		e.U16 = i;
		e.U32 = i;
		e.U64 = i;
		e.UGG = i;
		e.F32 = float(i);
		e.F64 = double(i);

		for (int ii = 0; ii != ARRAYSZ; ii++) {
			e.ArrayBs[ii]  = ii == i;
			e.ArrayI8[ii]  = -i;
			e.ArrayI16[ii] = -i;
			e.ArrayI32[ii] = -i;
			e.ArrayI64[ii] = -i;
			e.ArrayG64[ii] = -i;
			e.ArrayU8[ii]  = i;
			e.ArrayU16[ii] = i;
			e.ArrayU32[ii] = i;
			e.ArrayU64[ii] = i;
			e.ArrayUGG[ii] = i;
			e.ArrayF32[ii] = float(i);
			e.ArrayF64[ii] = double(i);
		}

		t->Fill();
	}

	f->Write();
	f->Close();

	exit(0);
}
`
